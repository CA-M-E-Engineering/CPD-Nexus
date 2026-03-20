package handlers

import (
	"encoding/json"
	"net/http"

	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"

	"github.com/gorilla/mux"
)

type ProjectsHandler struct {
	service ports.ProjectService
}

func NewProjectsHandler(service ports.ProjectService) *ProjectsHandler {
	return &ProjectsHandler{service: service}
}

func (h *ProjectsHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	userID := ports.GetUserID(r.Context())
	projects, err := h.service.ListProjects(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *ProjectsHandler) GetProjectById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID := ports.GetUserID(r.Context())

	project, err := h.service.GetProject(r.Context(), userID, id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectsHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project domain.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Enforce multi-tenancy: set UserID from authenticated context
	project.UserID = ports.GetUserID(r.Context())
	
	if err := h.service.CreateProject(r.Context(), &project); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectsHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID := ports.GetUserID(r.Context())

	var project domain.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateProject(r.Context(), userID, id, &project); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *ProjectsHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID := ports.GetUserID(r.Context())

	if err := h.service.DeleteProject(r.Context(), userID, id); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
