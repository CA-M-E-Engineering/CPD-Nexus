package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sgbuildex/internal/api/middleware"
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
	TimeIn  string `json:"time_in"`
	TimeOut string `json:"time_out"`
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

	// Resolve the owner user_id dynamically from context.
	// The bridge manager injects "bridge_userID" when dispatching messages.
	ownerID := ""
	if v, ok := ctx.Value("bridge_userID").(string); ok && v != "" {
		ownerID = v
	} else if v, ok := ctx.Value(middleware.UserIDKey).(string); ok && v != "" {
		ownerID = v
	}

	if ownerID == "" {
		log.Printf("AttendanceHandler: Cannot determine owner for attendance response, skipping")
		return nil, nil
	}

	// Elevate context to Vendor status for background database operations
	ctx = context.WithValue(ctx, middleware.UserIDKey, ownerID)
	ctx = context.WithValue(ctx, middleware.IsVendorKey, true)

	result := envelope.Content
	for _, rec := range result.Records {
		rawPayload, _ := json.Marshal(rec)

		err := h.service.ProcessBridgeAttendance(
			ctx,
			result.WorkerID,
			rec.TimeIn,
			rec.TimeOut,
			rawPayload,
		)

		if err != nil {
			log.Printf("AttendanceHandler: Service error for worker %s (owner %s): %v", result.WorkerID, ownerID, err)
		}
	}
	log.Printf("AttendanceHandler: Finished processing %d records for worker %s (owner: %s)", len(result.Records), result.WorkerID, ownerID)

	return nil, nil
}
