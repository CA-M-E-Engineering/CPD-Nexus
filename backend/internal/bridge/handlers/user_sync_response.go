package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"cpd-nexus/internal/bridge"
	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/pkg/logger"
	"strings"
)

// UserSyncResponsePayload maps the generic payload response from bridge
type UserSyncResponsePayload struct {
	Code    int             `json:"code"`
	Msg     string          `json:"msg"`
	Content json.RawMessage `json:"content"`
}

// UserSyncResponseHandler processes REGISTER_USER_RESPONSE and UPDATE_USER_RESPONSE
type UserSyncResponseHandler struct {
	workerRepo ports.WorkerRepository
	bridgeRepo ports.BridgeRepository
}

func NewUserSyncResponseHandler(workerRepo ports.WorkerRepository, bridgeRepo ports.BridgeRepository) *UserSyncResponseHandler {
	return &UserSyncResponseHandler{
		workerRepo: workerRepo,
		bridgeRepo: bridgeRepo,
	}
}

func (h *UserSyncResponseHandler) Handle(ctx context.Context, msg bridge.Message) (*bridge.Message, error) {
	var payload UserSyncResponsePayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user sync response: %w", err)
	}

	// RequestID is injected as "req-xxx|workerID" from BuildSyncRequests
	parts := strings.Split(msg.Meta.RequestID, "|")
	if len(parts) < 2 {
		logger.Infof("[UserSyncResponse] Warning: Cannot extract worker ID from request_id: %s", msg.Meta.RequestID)
		return nil, nil // Cannot determine which worker to update
	}

	workerID := parts[len(parts)-1] // Target worker ID

	if payload.Code == 200 {
		logger.Infof("[UserSyncResponse] Bridge returned success (200) for worker %s. Marking as synced.", workerID)

		// Hard update to mark worker as synced in DB
		if err := h.workerRepo.MarkSynced(ctx, workerID); err != nil {
			logger.Infof("[UserSyncResponse] Failed to update sync status for worker %s: %v", workerID, err)
			return nil, err
		}
	} else {
		// Do not update is_synced if bridge explicitly rejected or failed the user operation
		logger.Infof("[UserSyncResponse] Bridge rejected sync for worker %s (Code: %d, Msg: %s). Sync status unchanged.", workerID, payload.Code, payload.Msg)
	}

	// Update the interaction log with the actual status code
	// bridge_userID was injected into ctx by RequestManager
	if userID, ok := ctx.Value("bridge_userID").(string); ok {
		_ = h.bridgeRepo.LogBridgeInteraction(ctx, userID, msg.Action, msg.Meta.RequestID, nil, msg.Payload, payload.Code)
	}

	return nil, nil
}
