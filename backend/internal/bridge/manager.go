package bridge

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sgbuildex/internal/core/domain"
	"sync"
	"time"
)

// RequestManager handles business-level commands and response logic across multiple bridges
type RequestManager struct {
	Transports map[string]*Transport // Key: user_id
	DB         *sql.DB
	Handlers   map[string]Handler
	mu         sync.RWMutex
}

func NewRequestManager(db *sql.DB) *RequestManager {
	return &RequestManager{
		Transports: make(map[string]*Transport),
		DB:         db,
		Handlers:   make(map[string]Handler),
	}
}

// RegisterHandler adds a handler for a specific message type
func (rm *RequestManager) RegisterHandler(msgType string, h Handler) {
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

// RequestAttendance sends high-level commands to the appropriate bridge for each active worker
func (rm *RequestManager) RequestAttendance() error {
	log.Println("RequestManager: Starting attendance fetch for all workers across all bridges")

	query := `
		SELECT w.worker_id, w.user_id, p.site_id 
		FROM workers w
		JOIN projects p ON w.current_project_id = p.project_id
		WHERE w.status = 'active' AND w.current_project_id IS NOT NULL`

	rows, err := rm.DB.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query active workers: %w", err)
	}
	defer rows.Close()

	type workerTask struct {
		workerID string
		userID   string
		siteID   string
	}
	var tasks []workerTask
	for rows.Next() {
		var t workerTask
		if err := rows.Scan(&t.workerID, &t.userID, &t.siteID); err != nil {
			log.Printf("RequestManager: Error scanning worker: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}

	if len(tasks) == 0 {
		log.Println("RequestManager: No active workers with assigned projects found")
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
		transport, exists := rm.GetTransport(task.userID)
		if !exists || !transport.IsConnected() {
			log.Printf("RequestManager: Skipping worker %s, no active bridge connection for owner %s", task.workerID, task.userID)
			continue
		}

		// Get device SNs
		deviceRows, err := rm.DB.Query("SELECT sn FROM devices WHERE site_id = ? AND status != 'inactive'", task.siteID)
		if err != nil {
			continue
		}

		var sns []string
		for deviceRows.Next() {
			var sn string
			if err := deviceRows.Scan(&sn); err == nil {
				sns = append(sns, sn)
			}
		}
		deviceRows.Close()

		if len(sns) == 0 {
			continue
		}

		payload := map[string]interface{}{
			"worker_id":  task.workerID,
			"devices":    sns,
			"start_time": timeRange["from"],
			"end_time":   timeRange["to"],
		}

		req, err := NewRequest("GET_ATTENDANCE", payload)
		if err != nil {
			continue
		}

		if err := transport.Write(req); err != nil {
			log.Printf("RequestManager: Failed to send request for worker %s via bridge %s: %v", task.workerID, task.userID, err)
		} else {
			reqJSON, _ := json.MarshalIndent(req, "", "  ")
			log.Printf("\n--- [BRIDGE OUTBOUND REQUEST (%s)] ---\n%s\n---------------------------------", task.userID, string(reqJSON))
		}
	}

	return nil
}

// RequestUserSync sends REGISTER_USER and UPDATE_USER commands for pending workers
func (rm *RequestManager) RequestUserSync(ctx context.Context, builder interface {
	BuildSyncRequests(ctx context.Context, userID string) ([]Message, []string, []domain.Worker, []domain.Worker, error)
	MarkWorkersSynced(ctx context.Context, workerIDs []string)
}) error {
	log.Println("RequestManager: Starting user sync check")

	messages, workerIDs, invalidWorkers, _, err := builder.BuildSyncRequests(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to build user sync requests: %w", err)
	}

	if len(invalidWorkers) > 0 {
		log.Printf("RequestManager: %d workers are missing devices and cannot be synced", len(invalidWorkers))
	}

	if len(messages) == 0 {
		log.Println("RequestManager: No pending user sync requests")
		return nil
	}

	log.Printf("RequestManager: Sending %d user sync requests", len(messages))

	// Because BuildSyncRequests currently doesn't return the UserID per worker easily,
	// we will need to lookup the user_id for the worker before sending, or handle it via Transport mapping.
	// We'll write a quick query to group messages by worker's user_id.

	var successIDs []string

	for i, msg := range messages {
		workerID := workerIDs[i]

		var ownerID string
		err := rm.DB.QueryRow("SELECT user_id FROM workers WHERE worker_id = ?", workerID).Scan(&ownerID)
		if err != nil {
			log.Printf("RequestManager: Cannot find owner for worker %s: %v", workerID, err)
			continue
		}

		transport, exists := rm.GetTransport(ownerID)
		if !exists || !transport.IsConnected() {
			log.Printf("RequestManager: Skipping sync for worker %s, bridge %s is disconnected", workerID, ownerID)
			continue
		}

		if err := transport.Write(msg); err != nil {
			log.Printf("RequestManager: Failed to send user sync request for %s: %v", workerID, err)
		} else {
			respMsg, _ := json.MarshalIndent(msg, "", "  ")
			log.Printf("\n--- [BRIDGE OUTBOUND USER SYNC (%s)] ---\n%s\n-----------------------------------", ownerID, string(respMsg))
			successIDs = append(successIDs, workerID)
		}
	}

	if len(successIDs) > 0 {
		builder.MarkWorkersSynced(ctx, successIDs)
		log.Printf("RequestManager: Marked %d workers as synced", len(successIDs))
	}

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
				log.Printf("RequestManager (%s): Read error: %v", userID, err)
				transport.Close()
				continue
			}

			fullMsg, _ := json.MarshalIndent(msg, "", "  ")
			log.Printf("\n--- [BRIDGE INBOUND (%s)] ---\n%s\n------------------------", userID, string(fullMsg))

			if handler, ok := rm.Handlers[msg.Action]; ok {
				// Setting up an extended context to pass the owner ID to handlers if needed
				reqCtx := context.WithValue(ctx, "bridge_userID", userID)
				resp, err := handler.Handle(reqCtx, msg)
				if err != nil {
					log.Printf("RequestManager (%s): Handler for %s failed: %v", userID, msg.Action, err)
				} else if resp != nil {
					if err := transport.Write(*resp); err != nil {
						log.Printf("RequestManager (%s): Failed to send response back to bridge: %v", userID, err)
					} else {
						respMsg, _ := json.MarshalIndent(resp, "", "  ")
						log.Printf("\n--- [BRIDGE OUTBOUND RESPONSE (%s)] ---\n%s\n----------------------------------", userID, string(respMsg))
					}
				}
			} else {
				log.Printf("RequestManager (%s): Received unknown action: %s", userID, msg.Action)
			}
		}
	}
}
