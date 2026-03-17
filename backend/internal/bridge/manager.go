package bridge

import (
	"context"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/pkg/logger"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// RequestManager handles business-level commands and response logic across multiple bridges
type RequestManager struct {
	Transports map[string]*Transport // Key: user_id
	BridgeRepo ports.BridgeRepository
	Handlers   map[string]Handler
	mu         sync.RWMutex // protects Transports
	handlersMu sync.RWMutex // protects Handlers
}

func NewRequestManager(bridgeRepo ports.BridgeRepository) *RequestManager {
	return &RequestManager{
		Transports: make(map[string]*Transport),
		BridgeRepo: bridgeRepo,
		Handlers:   make(map[string]Handler),
	}
}

// RegisterHandler adds a handler for a specific message type (thread-safe)
func (rm *RequestManager) RegisterHandler(msgType string, h Handler) {
	rm.handlersMu.Lock()
	defer rm.handlersMu.Unlock()
	rm.Handlers[msgType] = h
}

// AddTransport adds a new transport for a user
func (rm *RequestManager) AddTransport(userID string, t *Transport) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Close existing if replacing
	if existing, ok := rm.Transports[userID]; ok {
		existing.Close()
	}
	rm.Transports[userID] = t
}

// RemoveTransport removes a transport for a user
func (rm *RequestManager) RemoveTransport(userID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if existing, ok := rm.Transports[userID]; ok {
		existing.Close()
		delete(rm.Transports, userID)
	}
}

// GetTransport gets a transport for a specific user safely
func (rm *RequestManager) GetTransport(userID string) (*Transport, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	t, ok := rm.Transports[userID]
	return t, ok
}

// GetAllTransports returns a copy of all transports for iteration
func (rm *RequestManager) GetAllTransports() map[string]*Transport {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	copyMap := make(map[string]*Transport)
	for k, v := range rm.Transports {
		copyMap[k] = v
	}
	return copyMap
}

// RequestAttendance sends high-level commands to the appropriate bridge for each active worker.
// ctx should be derived from the main application context to support graceful shutdown.
func (rm *RequestManager) RequestAttendance(ctx context.Context) error {
	logger.Infof("RequestManager: Starting attendance fetch for all workers across all bridges")

	tasks, err := rm.BridgeRepo.GetActiveBridgeWorkers(ctx)
	if err != nil {
		return fmt.Errorf("failed to query active workers: %w", err)
	}

	if len(tasks) == 0 {
		logger.Infof("RequestManager: No active workers with assigned projects found")
		return nil
	}

	// Group tasks by UserID to process bridges in parallel but workers sequentially per bridge
	tasksByOwner := make(map[string][]ports.BridgeWorkerTask)
	for _, task := range tasks {
		tasksByOwner[task.UserID] = append(tasksByOwner[task.UserID], task)
	}

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	startTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
	timeRangeFrom := startTime.Format(time.RFC3339)
	timeRangeTo := now.Format(time.RFC3339)

	var wg sync.WaitGroup
	for ownerID, ownerTasks := range tasksByOwner {
		wg.Add(1)
		go func(uid string, workerTasks []ports.BridgeWorkerTask) {
			defer wg.Done()

			transport, exists := rm.GetTransport(uid)
			if !exists || !transport.IsConnected() {
				logger.Infof("RequestManager (%s): Skipping fetch, bridge not connected", uid)
				return
			}

			for _, task := range workerTasks {
				// Check for cancellation between workers
				select {
				case <-ctx.Done():
					return
				default:
				}

				sns, err := rm.BridgeRepo.GetActiveDeviceSNsBySite(ctx, task.SiteID)
				if err != nil || len(sns) == 0 {
					continue
				}

				payload := map[string]interface{}{
					"worker_id":  task.WorkerID,
					"devices":    sns,
					"start_time": timeRangeFrom,
					"end_time":   timeRangeTo,
				}

				req, err := NewRequest("GET_ATTENDANCE", payload)
				if err != nil {
					continue
				}

				if err := transport.Write(req); err != nil {
					logger.Infof("RequestManager (%s): Failed to send request for worker %s: %v", uid, task.WorkerID, err)
				} else {
					logger.Infof("RequestManager (%s): Queued attendance request for worker %s", uid, task.WorkerID)
				}
			}
		}(ownerID, ownerTasks)
	}

	wg.Wait()
	return nil
}

// RequestUserSync sends REGISTER_USER and UPDATE_USER commands for pending workers
func (rm *RequestManager) RequestUserSync(ctx context.Context, builder interface {
	BuildSyncRequests(ctx context.Context, userID string) ([]Message, []string, []domain.Worker, []domain.Worker, error)
}) error {
	logger.Infof("RequestManager: Starting user sync check")

	messages, workerIDs, invalidWorkers, _, err := builder.BuildSyncRequests(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to build user sync requests: %w", err)
	}

	if len(invalidWorkers) > 0 {
		logger.Infof("RequestManager: %d workers are missing devices and cannot be synced", len(invalidWorkers))
	}

	if len(messages) == 0 {
		logger.Infof("RequestManager: No pending user sync requests")
		return nil
	}

	// Group messages by owner for parallel delivery across bridges
	type syncTask struct {
		workerID string
		msg      Message
	}
	tasksByOwner := make(map[string][]syncTask)

	for i, msg := range messages {
		workerID := workerIDs[i]
		ownerID, err := rm.BridgeRepo.GetWorkerOwnerID(ctx, workerID)
		if err != nil {
			logger.Infof("RequestManager: Cannot find owner for worker %s: %v", workerID, err)
			continue
		}
		tasksByOwner[ownerID] = append(tasksByOwner[ownerID], syncTask{workerID: workerID, msg: msg})
	}

	var wg sync.WaitGroup
	for ownerID, tasks := range tasksByOwner {
		wg.Add(1)
		go func(uid string, subTasks []syncTask) {
			defer wg.Done()

			transport, exists := rm.GetTransport(uid)
			if !exists || !transport.IsConnected() {
				logger.Infof("RequestManager (%s): Skipping user sync, bridge not connected", uid)
				return
			}

			for _, t := range subTasks {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if err := transport.Write(t.msg); err != nil {
					logger.Infof("RequestManager (%s): Failed to send sync for worker %s: %v", uid, t.workerID, err)
				} else {
					logger.Infof("RequestManager (%s): Queued sync for worker %s", uid, t.workerID)
				}
			}
		}(ownerID, tasks)
	}

	wg.Wait()
	return nil
}

// HandleIncomingMessages needs to be updated. Now we'll run a listener per transport.
func (rm *RequestManager) HandleIncomingMessages(ctx context.Context, userID string, transport *Transport) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !transport.IsConnected() {
				time.Sleep(2 * time.Second)
				continue
			}

			msg, err := transport.Read()
			if err != nil {
				logger.Infof("RequestManager (%s): Read error: %v", userID, err)
				transport.Close()
				continue
			}

			fullMsg, _ := json.MarshalIndent(msg, "", "  ")
			logger.Infof("\n--- [BRIDGE INBOUND (%s)] ---\n%s\n------------------------", userID, string(fullMsg))

			if handler, ok := func() (Handler, bool) {
				rm.handlersMu.RLock()
				defer rm.handlersMu.RUnlock()
				h, ok := rm.Handlers[msg.Action]
				return h, ok
			}(); ok {
				// Setting up an extended context to pass the owner ID to handlers if needed
				reqCtx := context.WithValue(ctx, "bridge_userID", userID)
				resp, err := handler.Handle(reqCtx, msg)
				if err != nil {
					logger.Infof("RequestManager (%s): Handler for %s failed: %v", userID, msg.Action, err)
				} else if resp != nil {
					if err := transport.Write(*resp); err != nil {
						logger.Infof("RequestManager (%s): Failed to send response back to bridge: %v", userID, err)
					} else {
						respMsg, _ := json.MarshalIndent(resp, "", "  ")
						logger.Infof("\n--- [BRIDGE OUTBOUND RESPONSE (%s)] ---\n%s\n----------------------------------", userID, string(respMsg))
					}
				}
			} else {
				logger.Infof("RequestManager (%s): Received unknown action: %s", userID, msg.Action)
			}
		}
	}
}
