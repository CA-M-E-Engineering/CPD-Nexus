package handlers

import (
	"encoding/json"
	"net/http"

	"cpd-nexus/internal/core/ports"

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
	userID := ports.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized: valid user context required", http.StatusUnauthorized)
		return
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

	userID := ports.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusUnauthorized)
		return
	}

	// Admin/Vendor bypass: if they are an admin, they can test any project
	if ports.IsVendor(r.Context()) {
		userID = ""
	}

	submitted, failed, err := h.pitstopService.TestSubmission(r.Context(), userID, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Manual submission complete",
		"metrics": map[string]int{
			"payloads_submitted": submitted,
			"validation_failed":  failed,
		},
	})
}

// GetTestingProjects handles retrieving a list of unique projects that currently have pending attendance records
func (h *PitstopHandler) GetTestingProjects(w http.ResponseWriter, r *http.Request) {
	userID := ports.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusUnauthorized)
		return
	}

	// Admin/Vendor bypass: if they are an admin, they can see all projects
	if ports.IsVendor(r.Context()) {
		userID = ""
	}

	projects, err := h.pitstopService.GetProjectsWithPendingAttendance(r.Context(), userID)
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
