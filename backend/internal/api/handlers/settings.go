package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type SettingsHandler struct {
	service ports.SettingsService
}

func NewSettingsHandler(service ports.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.service.GetSettings(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (h *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var payload domain.SystemSettings
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Enforce singleton ID constraint — this table always has a single row with ID=1.
	payload.ID = 1

	if err := h.service.UpdateSettings(r.Context(), payload); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}
