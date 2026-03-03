package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/core/services"
)

type PitstopHandler struct {
	pitstopService *services.PitstopService
}

func NewPitstopHandler(pitstopService *services.PitstopService) *PitstopHandler {
	return &PitstopHandler{
		pitstopService: pitstopService,
	}
}

// GetAuthorisations handles retrieving existing configuration mappings
func (h *PitstopHandler) GetAuthorisations(w http.ResponseWriter, r *http.Request) {
	auths, err := h.pitstopService.GetAuthorisations(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auths)
}

// SyncConfig handles triggering an immediate sync over external pitstop endpoints
func (h *PitstopHandler) SyncConfig(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "Owner_001"
	}

	if err := h.pitstopService.SyncConfig(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Pitstop configurations synced successfully"}`))
}
