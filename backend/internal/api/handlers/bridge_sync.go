package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sgbuildex/internal/bridge"
	bridgeHandlers "sgbuildex/internal/bridge/handlers"
)

// BridgeSyncHandler handles manual sync trigger from the frontend
type BridgeSyncHandler struct {
	builder   *bridgeHandlers.UserSyncBuilder
	transport *bridge.Transport
}

func NewBridgeSyncHandler(builder *bridgeHandlers.UserSyncBuilder, transport *bridge.Transport) *BridgeSyncHandler {
	return &BridgeSyncHandler{
		builder:   builder,
		transport: transport,
	}
}

// SyncUsers handles POST /api/bridge/sync-users
func (h *BridgeSyncHandler) SyncUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")

	// Get UserID from header (passed by frontend)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Missing X-User-ID header",
		})
		return
	}

	// Check if bridge is connected
	if !h.transport.IsConnected() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Bridge is not connected",
		})
		return
	}

	// Build sync requests with userID filter
	messages, workerIDs, invalidWorkers, unauthWorkers, err := h.builder.BuildSyncRequests(ctx, userID)
	if err != nil {
		log.Printf("[BridgeSync API] Failed to build sync requests: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Logging issues but NOT aborting
	if len(invalidWorkers) > 0 || len(unauthWorkers) > 0 {
		log.Printf("[BridgeSync API] Found issues: %d missing devices, %d missing auth", len(invalidWorkers), len(unauthWorkers))
	}

	if len(messages) == 0 && len(invalidWorkers) == 0 && len(unauthWorkers) == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"message":  "No workers pending sync.",
			"sent":     0,
			"register": 0,
			"update":   0,
		})
		return
	}

	// Send each message via bridge transport
	var successIDs []string
	registerCount := 0
	updateCount := 0
	failCount := 0

	for i, msg := range messages {
		if err := h.transport.Write(msg); err != nil {
			log.Printf("[BridgeSync API] Failed to send message %d: %v", i, err)
			failCount++
		} else {
			respMsg, _ := json.MarshalIndent(msg, "", "  ")
			log.Printf("\n--- [BRIDGE SYNC API OUTBOUND] ---\n%s\n----------------------------------", string(respMsg))

			if i < len(workerIDs) {
				successIDs = append(successIDs, workerIDs[i])
			}

			if msg.Action == "REGISTER_USER" {
				registerCount++
			} else {
				updateCount++
			}
		}
	}

	// Mark successfully sent workers as synced (is_synced=1)
	if len(successIDs) > 0 {
		h.builder.MarkWorkersSynced(ctx, successIDs)
		log.Printf("[BridgeSync API] Marked %d workers as synced", len(successIDs))
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":                 true,
		"message":                 "Sync completed",
		"sent":                    len(successIDs),
		"register":                registerCount,
		"update":                  updateCount,
		"failed":                  failCount,
		"invalid_workers":         invalidWorkers,
		"unauthenticated_workers": unauthWorkers,
	})
}
