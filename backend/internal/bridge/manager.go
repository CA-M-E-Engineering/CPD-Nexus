package bridge

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sgbuildex/internal/core/domain"
	"time"
)

// RequestManager handles business-level commands and response logic
type RequestManager struct {
	Transport *Transport
	DB        *sql.DB
	Handlers  map[string]Handler
}

func NewRequestManager(t *Transport, db *sql.DB) *RequestManager {
	return &RequestManager{
		Transport: t,
		DB:        db,
		Handlers:  make(map[string]Handler),
	}
}

// RegisterHandler adds a handler for a specific message type
func (rm *RequestManager) RegisterHandler(msgType string, h Handler) {
	rm.Handlers[msgType] = h
}

// RequestAttendance sends high-level commands to the bridge for each active worker
func (rm *RequestManager) RequestAttendance() error {
	log.Println("RequestManager: Starting attendance fetch for all workers")

	// 1. Get all active workers
	// We use the DB directly or a repo if passed in. Since we only have rm.DB, we query directly.
	query := `
		SELECT w.worker_id, p.site_id 
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
		siteID   string
	}
	var tasks []workerTask
	for rows.Next() {
		var t workerTask
		if err := rows.Scan(&t.workerID, &t.siteID); err != nil {
			log.Printf("RequestManager: Error scanning worker: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}

	if len(tasks) == 0 {
		log.Println("RequestManager: No active workers with assigned projects found")
		return nil
	}

	// 2. Build time range for today
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	timeRange := map[string]string{
		"from": midnight.Format(time.RFC3339),
		"to":   now.Format(time.RFC3339),
	}

	// 3. For each worker, get their site's devices and send a request
	for _, task := range tasks {
		// Get device SNs for this site
		deviceRows, err := rm.DB.Query("SELECT sn FROM devices WHERE site_id = ? AND status != 'inactive'", task.siteID)
		if err != nil {
			log.Printf("RequestManager: Failed to query devices for site %s: %v", task.siteID, err)
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
			log.Printf("RequestManager: No devices found for worker %s at site %s", task.workerID, task.siteID)
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
			log.Printf("RequestManager: Failed to create request for worker %s: %v", task.workerID, err)
			continue
		}

		if err := rm.Transport.Write(req); err != nil {
			log.Printf("RequestManager: Failed to send request for worker %s: %v", task.workerID, err)
		} else {
			reqJSON, _ := json.MarshalIndent(req, "", "  ")
			log.Printf("\n--- [BRIDGE OUTBOUND REQUEST] ---\n%s\n---------------------------------", string(reqJSON))
			log.Printf("RequestManager: Sent GET_ATTENDANCE request for worker %s on %d devices", task.workerID, len(sns))
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

	// For background sync, we pass empty userID to sync all
	messages, workerIDs, invalidWorkers, _, err := builder.BuildSyncRequests(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to build user sync requests: %w", err)
	}

	if len(invalidWorkers) > 0 {
		log.Printf("RequestManager: %d workers are missing devices and cannot be synced", len(invalidWorkers))
		for _, w := range invalidWorkers {
			log.Printf("  - Worker %s (%s) at site %s", w.ID, w.Name, w.SiteID)
		}
	}

	if len(messages) == 0 {
		log.Println("RequestManager: No pending user sync requests")
		return nil
	}

	log.Printf("RequestManager: Sending %d user sync requests", len(messages))

	var successIDs []string
	for i, msg := range messages {
		if err := rm.Transport.Write(msg); err != nil {
			log.Printf("RequestManager: Failed to send user sync request %d: %v", i, err)
		} else {
			respMsg, _ := json.MarshalIndent(msg, "", "  ")
			log.Printf("\n--- [BRIDGE OUTBOUND USER SYNC] ---\n%s\n-----------------------------------", string(respMsg))
			if i < len(workerIDs) {
				successIDs = append(successIDs, workerIDs[i])
			}
		}
	}

	// Mark successfully sent workers as synced
	if len(successIDs) > 0 {
		builder.MarkWorkersSynced(ctx, successIDs)
		log.Printf("RequestManager: Marked %d workers as synced", len(successIDs))
	}

	return nil
}

// HandleIncomingMessages dispatches received messages to their respective handlers
func (rm *RequestManager) HandleIncomingMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !rm.Transport.IsConnected() {
				time.Sleep(2 * time.Second)
				continue
			}

			msg, err := rm.Transport.Read()
			if err != nil {
				log.Printf("RequestManager: Read error: %v. Transport may handle reconnection.", err)
				rm.Transport.Close()
				continue
			}

			fullMsg, _ := json.MarshalIndent(msg, "", "  ")
			log.Printf("\n--- [BRIDGE INBOUND] ---\n%s\n------------------------", string(fullMsg))

			if handler, ok := rm.Handlers[msg.Action]; ok {
				resp, err := handler.Handle(ctx, msg)
				if err != nil {
					log.Printf("RequestManager: Handler for %s failed: %v", msg.Action, err)
				} else if resp != nil {
					// Send response back if handler provided one
					if err := rm.Transport.Write(*resp); err != nil {
						log.Printf("RequestManager: Failed to send response back to bridge: %v", err)
					} else {
						respMsg, _ := json.MarshalIndent(resp, "", "  ")
						log.Printf("\n--- [BRIDGE OUTBOUND RESPONSE] ---\n%s\n----------------------------------", string(respMsg))
					}
				}
			} else {
				log.Printf("RequestManager: Received unknown action: %s", msg.Action)
			}
		}
	}
}
