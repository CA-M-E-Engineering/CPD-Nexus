package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/api/middleware"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/apperrors"

	"github.com/gorilla/mux"
)

type ProjectsHandler struct {
	service ports.ProjectService
}

func NewProjectsHandler(service ports.ProjectService) *ProjectsHandler {
	return &ProjectsHandler{service: service}
}

func (h *ProjectsHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	projects, err := h.service.ListProjects(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *ProjectsHandler) GetProjectById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID := middleware.GetUserID(r.Context())
	project, err := h.service.GetProject(r.Context(), userID, id)
	if err != nil {
		h.handleError(w, err)
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

	if err := h.service.CreateProject(r.Context(), &project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectsHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var project domain.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r.Context())
	if err := h.service.UpdateProject(r.Context(), userID, id, &project); err != nil {
		h.handleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *ProjectsHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID := middleware.GetUserID(r.Context())
	if err := h.service.DeleteProject(r.Context(), userID, id); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
func (h *ProjectsHandler) handleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
