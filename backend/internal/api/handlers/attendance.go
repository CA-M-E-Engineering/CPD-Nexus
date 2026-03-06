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
	// userID comes from the middleware context, not a query param, for tenant isolation.
	userID := ports.GetUserID(r.Context())
	siteID := r.URL.Query().Get("site_id")
	workerID := r.URL.Query().Get("worker_id")
	date := r.URL.Query().Get("date")

	records, err := h.service.ListAttendance(r.Context(), userID, siteID, workerID, date)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}
