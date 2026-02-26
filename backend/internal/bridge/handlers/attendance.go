package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sgbuildex/internal/bridge"
	"sgbuildex/internal/core/ports"
)

// AttendanceResult matches the top-level payload from the bridge
type AttendanceResult struct {
	WorkerID string            `json:"worker_id"`
	Devices  []string          `json:"devices"`
	Records  []AttendanceEvent `json:"records"`
}

// AttendanceEvent matches the individual record structure from the bridge
type AttendanceEvent struct {
	DeviceID string `json:"device_id"`
	PersonID string `json:"person_id"`
	TimeIn   string `json:"time_in"`
	TimeOut  string `json:"time_out"`
}

type AttendanceHandler struct {
	service ports.AttendanceService
}

func NewAttendanceHandler(service ports.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) Handle(ctx context.Context, msg bridge.Message) (*bridge.Message, error) {
	// The bridge returns a wrapped response: { "code": 200, "msg": "...", "content": AttendanceResult }
	var envelope struct {
		Code    int              `json:"code"`
		Msg     string           `json:"msg"`
		Content AttendanceResult `json:"content"`
	}

	if err := json.Unmarshal(msg.Payload, &envelope); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attendance envelope: %w", err)
	}

	if envelope.Code != 200 {
		log.Printf("AttendanceHandler: Bridge returned error %d: %s", envelope.Code, envelope.Msg)
		return nil, nil
	}

	result := envelope.Content
	// Since the backend sends one request per worker and the bridge filters by worker_id,
	// we can process these records directly.
	for _, rec := range result.Records {
		// Ensure we only process records for the requested worker (extra safety)
		if result.WorkerID != "" && rec.PersonID != result.WorkerID {
			continue
		}

		// Use the service to process and persist the attendance record
		rawPayload, _ := json.Marshal(rec)

		err := h.service.ProcessBridgeAttendance(
			ctx,
			rec.DeviceID,
			rec.PersonID,
			rec.TimeIn,
			rec.TimeOut,
			rawPayload,
		)

		if err != nil {
			log.Printf("AttendanceHandler: Service error for worker %s on device %s: %v", rec.PersonID, rec.DeviceID, err)
		}
	}
	log.Printf("AttendanceHandler: Finished processing %d records for worker %s", len(result.Records), result.WorkerID)

	return nil, nil
}
