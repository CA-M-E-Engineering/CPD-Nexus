package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/logger"
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

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	startTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())

	timeRange := map[string]string{
		"from": startTime.Format(time.RFC3339),
		"to":   now.Format(time.RFC3339),
	}

	// Create requests
	for _, task := range tasks {
		transport, exists := rm.GetTransport(task.UserID)
		if !exists || !transport.IsConnected() {
			logger.Infof("RequestManager: Skipping worker %s, no active bridge connection for owner %s", task.WorkerID, task.UserID)
			continue
		}

		// Get device SNs
		sns, err := rm.BridgeRepo.GetActiveDeviceSNsBySite(ctx, task.SiteID)
		if err != nil {
			continue
		}

		if len(sns) == 0 {
			continue
		}

		payload := map[string]interface{}{
			"worker_id":  task.WorkerID,
			"devices":    sns,
			"start_time": timeRange["from"],
			"end_time":   timeRange["to"],
		}

		req, err := NewRequest("GET_ATTENDANCE", payload)
		if err != nil {
			continue
		}

		if err := transport.Write(req); err != nil {
			logger.Infof("RequestManager: Failed to send request for worker %s via bridge %s: %v", task.WorkerID, task.UserID, err)
		} else {
			reqJSON, _ := json.MarshalIndent(req, "", "  ")
			logger.Infof("\n--- [BRIDGE OUTBOUND REQUEST (%s)] ---\n%s\n---------------------------------", task.UserID, string(reqJSON))
		}
	}

	return nil
}

// RequestUserSync sends REGISTER_USER and UPDATE_USER commands for pending workers
func (rm *RequestManager) RequestUserSync(ctx context.Context, builder interface {
	BuildSyncRequests(ctx context.Context, userID string) ([]Message, []string, []domain.Worker, []domain.Worker, error)
	MarkWorkersSynced(ctx context.Context, workerIDs []string)
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

	logger.Infof("RequestManager: Sending %d user sync requests", len(messages))

	// Because BuildSyncRequests currently doesn't return the UserID per worker easily,
	// we will need to lookup the user_id for the worker before sending, or handle it via Transport mapping.
	// We'll write a quick query to group messages by worker's user_id.

	var successIDs []string

	for i, msg := range messages {
		workerID := workerIDs[i]

		ownerID, err := rm.BridgeRepo.GetWorkerOwnerID(ctx, workerID)
		if err != nil {
			logger.Infof("RequestManager: Cannot find owner for worker %s: %v", workerID, err)
			continue
		}

		transport, exists := rm.GetTransport(ownerID)
		if !exists || !transport.IsConnected() {
			logger.Infof("RequestManager: Skipping sync for worker %s, bridge %s is disconnected", workerID, ownerID)
			continue
		}

		if err := transport.Write(msg); err != nil {
			logger.Infof("RequestManager: Failed to send user sync request for %s: %v", workerID, err)
		} else {
			respMsg, _ := json.MarshalIndent(msg, "", "  ")
			logger.Infof("\n--- [BRIDGE OUTBOUND USER SYNC (%s)] ---\n%s\n-----------------------------------", ownerID, string(respMsg))
			successIDs = append(successIDs, workerID)
		}
	}

	logger.Infof("RequestManager: Successfully queued %d worker sync requests across all bridges", len(successIDs))
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
