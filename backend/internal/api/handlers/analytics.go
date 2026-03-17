package handlers

import (
	"encoding/json"
	"net/http"
	"cpd-nexus/internal/core/ports"
)

type AnalyticsHandler struct {
	service ports.AnalyticsService
}

func NewAnalyticsHandler(service ports.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	stats, err := h.service.GetDashboardStats(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *AnalyticsHandler) GetActivityLog(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	isVendor := ports.IsVendor(r.Context())
	contextUserID := ports.GetUserID(r.Context())

	// Security: If not vendor, force filter to context user ID regardless of what query says
	if !isVendor {
		userID = contextUserID
	}

	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	filters := make(map[string]interface{})
	if action := r.URL.Query().Get("action"); action != "" {
		filters["action"] = action
	}
	if targetType := r.URL.Query().Get("target_type"); targetType != "" {
		filters["target_type"] = targetType
	}

	logs, err := h.service.GetActivityLog(r.Context(), userID, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *AnalyticsHandler) GetDetailedAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetDetailedAnalytics(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
