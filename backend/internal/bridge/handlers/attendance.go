package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sgbuildex/internal/bridge"
	"sgbuildex/internal/core/ports"
	"time"
)

// AttendanceResult matches the top-level payload from the bridge
type AttendanceResult struct {
	DeviceID string            `json:"device_id"`
	Records  []AttendanceEvent `json:"records"`
}

// AttendanceEvent matches the individual record structure from the bridge
type AttendanceEvent struct {
	PersonID string    `json:"person_id"`
	TimeIn   time.Time `json:"time_in"`
	TimeOut  time.Time `json:"time_out"`
}

type AttendanceHandler struct {
	service ports.AttendanceService
}

func NewAttendanceHandler(service ports.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) Handle(ctx context.Context, msg bridge.Message) (*bridge.Message, error) {
	var result AttendanceResult
	if err := json.Unmarshal(msg.Payload, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attendance result: %w", err)
	}

	for _, rec := range result.Records {
		// Use the service to process and persist the attendance record
		rawPayload, _ := json.Marshal(rec)

		err := h.service.ProcessBridgeAttendance(
			ctx,
			result.DeviceID,
			rec.PersonID,
			rec.TimeIn.Format(time.RFC3339),
			rec.TimeOut.Format(time.RFC3339),
			rawPayload,
		)

		if err != nil {
			log.Printf("AttendanceHandler: Service error for %s: %v", rec.PersonID, err)
		}
	}
	log.Printf("AttendanceHandler: Finished processing %d records for device %s", len(result.Records), result.DeviceID)

	return nil, nil
}
