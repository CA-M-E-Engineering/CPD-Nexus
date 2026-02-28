package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sgbuildex/internal/api/middleware"
	"sgbuildex/internal/bridge"
	bridgeHandlers "sgbuildex/internal/bridge/handlers"
)

// BridgeSyncHandler handles manual sync trigger from the frontend
type BridgeSyncHandler struct {
	builder    *bridgeHandlers.UserSyncBuilder
	requestMgr *bridge.RequestManager
}

func NewBridgeSyncHandler(builder *bridgeHandlers.UserSyncBuilder, requestMgr *bridge.RequestManager) *BridgeSyncHandler {
	return &BridgeSyncHandler{
		builder:    builder,
		requestMgr: requestMgr,
	}
}

// SyncUsers handles POST /api/bridge/sync-users
func (h *BridgeSyncHandler) SyncUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	// Use middleware to get UserID
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Missing User-ID scope",
		})
		return
	}

	// Check if bridge is connected for this user
	transport, exists := h.requestMgr.GetTransport(userID)
	if !exists || !transport.IsConnected() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Bridge is not connected for this account",
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
		if err := transport.Write(msg); err != nil {
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

	// Note: We no longer mark workers as synced here.
	// That responsibility is now deferred to the async bridge response handlers,
	// which will only mark them synced upon receiving a 200 OK from the hardware device.

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
