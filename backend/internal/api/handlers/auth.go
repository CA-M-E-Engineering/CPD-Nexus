package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/core/ports"
)

type AuthHandler struct {
	service ports.AuthService
}

func NewAuthHandler(service ports.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// LoginRequest structure
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse structure
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, user, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Simple role mapping logic for frontend legacy compatibility
	userMap := map[string]interface{}{
		"id":        user.ID,
		"tenant_id": user.ID,
		"name":      user.Name,
		"username":  user.Username,
		"role":      "manager", // Default
	}
	if user.TenantType == "client" {
		userMap["role"] = "client"
	}

	response := LoginResponse{
		Token: token,
		User:  userMap,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// For MVP: Return a static/mock user since we are not parsing JWT middleware yet
	user := map[string]interface{}{
		"id":        "testt.ltd",
		"tenant_id": "testt.ltd",
		"name":      "Test Tenant Admin",
		"username":  "testt.ltd",
		"role":      "client",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
