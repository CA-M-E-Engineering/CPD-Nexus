package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type SettingsHandler struct {
	Service ports.SettingsService
}

func NewSettingsHandler(service ports.SettingsService) *SettingsHandler {
	return &SettingsHandler{Service: service}
}

func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.Service.GetSettings(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Ensure ID is 1 (singleton)
	payload.ID = 1

	if err := h.Service.UpdateSettings(r.Context(), payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}
