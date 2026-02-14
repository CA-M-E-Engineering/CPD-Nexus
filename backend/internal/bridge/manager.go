package bridge

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

// RequestAttendance sends high-level commands to the bridge for each registered device
func (rm *RequestManager) RequestAttendance() error {
	log.Println("RequestManager: Starting attendance fetch for all devices")

	// 1. Get all device SNs from database
	rows, err := rm.DB.Query("SELECT sn FROM devices")
	if err != nil {
		return fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()

	var sns []string
	for rows.Next() {
		var sn string
		if err := rows.Scan(&sn); err != nil {
			log.Printf("RequestManager: Error scanning device SN: %v", err)
			continue
		}
		sns = append(sns, sn)
	}

	if len(sns) == 0 {
		log.Println("RequestManager: No devices found in database")
		return nil
	}

	// 2. Build time range for today
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	timeRange := map[string]string{
		"from": midnight.Format(time.RFC3339),
		"to":   now.Format(time.RFC3339),
	}

	// 3. Send a request for each device
	for _, sn := range sns {
		payload := map[string]interface{}{
			"device_id":  sn,
			"time_range": timeRange,
		}

		req, err := NewRequest("FETCH_ATTENDANCE", payload)
		if err != nil {
			log.Printf("RequestManager: Failed to create request for device %s: %v", sn, err)
			continue
		}

		if err := rm.Transport.Write(req); err != nil {
			log.Printf("RequestManager: Failed to send request for device %s: %v", sn, err)
		} else {
			reqMsg, _ := json.MarshalIndent(req, "", "  ")
			log.Printf("\n--- [BRIDGE OUTBOUND REQUEST] ---\n%s\n---------------------------------", string(reqMsg))
		}
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
