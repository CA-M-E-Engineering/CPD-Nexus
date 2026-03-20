package handlers

import (
	"context"
	"net/http"
	"cpd-nexus/internal/bridge"
	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/pkg/logger"

	"github.com/gorilla/websocket"
)

type BridgeHandler struct {
	requestMgr *bridge.RequestManager
	userRepo   ports.UserRepository
	upgrader   websocket.Upgrader
}

func NewBridgeHandler(mgr *bridge.RequestManager, userRepo ports.UserRepository) *BridgeHandler {
	return &BridgeHandler{
		requestMgr: mgr,
		userRepo:   userRepo,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // Allow all for now, can be tightened
		},
	}
}

// Connect handles the WebSocket upgrade from a local bridge
func (h *BridgeHandler) Connect(w http.ResponseWriter, r *http.Request) {
	logger.Infof("Bridge: Received connection request from %s", r.RemoteAddr)
	userID := r.URL.Query().Get("user_id")
	token := r.URL.Query().Get("token")

	if userID == "" || token == "" {
		http.Error(w, "Missing user_id or token", http.StatusBadRequest)
		return
	}

	// 1. Authenticate the bridge against the user record
	user, err := h.userRepo.Get(r.Context(), userID)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized: user not found", http.StatusUnauthorized)
		return
	}

	if user.BridgeAuthToken == nil || *user.BridgeAuthToken != token {
		http.Error(w, "Unauthorized: invalid bridge token", http.StatusUnauthorized)
		return
	}

	// 2. Upgrade to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("Bridge: Failed to upgrade connection for %s: %v", userID, err)
		return
	}

	logger.Infof("Bridge: New connection established for user %s", userID)

	// 3. Register transport in the manager
	t := bridge.NewServerTransport(conn, token)
	h.requestMgr.AddTransport(userID, t)

	// 4. Start message processing in the background
	// We use context.Background() here because the connection should live beyond the HTTP request lifecycle
	go h.requestMgr.HandleIncomingMessages(context.Background(), userID, t)
}
