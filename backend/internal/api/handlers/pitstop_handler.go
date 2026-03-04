package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/core/services"

	"github.com/gorilla/mux"
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

// AssignOnBehalfOf handles setting the user context for certain on_behalf_of names in pitstop configuration
func (h *PitstopHandler) AssignOnBehalfOf(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var input struct {
		OnBehalfOfNames []string `json:"on_behalf_of_names"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.pitstopService.AssignOnBehalfOfToUser(r.Context(), userID, input.OnBehalfOfNames); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "On behalf of entities successfully assigned to user"}`))
}
