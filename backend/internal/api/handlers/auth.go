package handlers

import (
	"encoding/json"
	"net/http"

	"cpd-nexus/internal/core/ports"
)

type AuthHandler struct {
	authService ports.AuthService
	userService ports.UserService
}

func NewAuthHandler(authService ports.AuthService, userService ports.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
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

	token, user, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Simple role mapping logic for frontend legacy compatibility
	userMap := map[string]interface{}{
		"id":       user.ID,
		"user_id":  user.ID,
		"name":     user.Name,
		"username": user.Username,
		"role":     "manager", // Default
	}
	switch user.UserType {
	case "client":
		userMap["role"] = "client"
	default:
		userMap["role"] = "manager"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	response := LoginResponse{
		Token: token, // keep for backward compatibility temporarily
		User:  userMap,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := ports.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Simple role mapping logic for frontend legacy compatibility
	userMap := map[string]interface{}{
		"id":       user.ID,
		"user_id":  user.ID,
		"name":     user.Name,
		"username": user.Username,
		"role":     "manager", // Default
	}
	switch user.UserType {
	case "client":
		userMap["role"] = "client"
	default:
		userMap["role"] = "manager"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userMap)
}
