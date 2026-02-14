package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/core/ports"
)

type AttendanceHandler struct {
	service ports.AttendanceService
}

func NewAttendanceHandler(service ports.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) GetAttendance(w http.ResponseWriter, r *http.Request) {
	// Query params
	siteID := r.URL.Query().Get("site_id")
	workerID := r.URL.Query().Get("worker_id")
	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	records, err := h.service.ListAttendance(r.Context(), userID, siteID, workerID, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}
