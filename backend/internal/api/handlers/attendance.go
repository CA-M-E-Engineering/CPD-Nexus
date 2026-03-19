package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/pkg/apperrors"
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

func (h *AttendanceHandler) UpdateAttendance(w http.ResponseWriter, r *http.Request) {
	userID := ports.GetUserID(r.Context())
	vars := mux.Vars(r)
	id := vars["id"]

	var payload struct {
		TimeIn  *time.Time `json:"time_in"`
		TimeOut *time.Time `json:"time_out"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, apperrors.NewValidationError("invalid request payload"))
		return
	}

	err := h.service.UpdateAttendance(r.Context(), userID, id, payload.TimeIn, payload.TimeOut)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
