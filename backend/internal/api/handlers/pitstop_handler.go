package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/core/ports"

	"github.com/gorilla/mux"
)

// PitstopHandler handles HTTP requests for Pitstop/SGBuildex operations.
// It depends on the ports.PitstopService interface, not the concrete service type.
type PitstopHandler struct {
	pitstopService ports.PitstopService
}

func NewPitstopHandler(pitstopService ports.PitstopService) *PitstopHandler {
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

// SyncConfig triggers an immediate sync with the external Pitstop configuration endpoint
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

// TestSubmission handles manual triggering of the CPD submission for a specific project
func (h *PitstopHandler) TestSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project_id"]

	if projectID == "" {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}

	count, err := h.pitstopService.TestSubmission(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Test submission completed successfully.",
		"metrics": map[string]int{
			"payloads_submitted": count,
		},
	}
	json.NewEncoder(w).Encode(response)
}

// GetTestingProjects handles retrieving a list of unique projects that currently have pending attendance records
func (h *PitstopHandler) GetTestingProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.pitstopService.GetProjectsWithPendingAttendance(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(projects) == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"data": projects})
}
