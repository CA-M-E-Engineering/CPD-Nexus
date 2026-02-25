package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/core/ports"

	"github.com/gorilla/mux"
)

type AssignmentsHandler struct {
	workerService  ports.WorkerService
	deviceService  ports.DeviceService
	projectService ports.ProjectService
}

func NewAssignmentsHandler(ws ports.WorkerService, ds ports.DeviceService, ps ports.ProjectService) *AssignmentsHandler {
	return &AssignmentsHandler{
		workerService:  ws,
		deviceService:  ds,
		projectService: ps,
	}
}

type AssignRequest struct {
	WorkerIDs  []string `json:"workerIds"`
	DeviceIDs  []string `json:"deviceIds"`
	ProjectIDs []string `json:"projectIds"`
}

func (h *AssignmentsHandler) AssignWorkers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["projectId"]

	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.workerService.AssignWorkersToProject(r.Context(), projectId, req.WorkerIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "assigned", "project": projectId})
}

func (h *AssignmentsHandler) AssignDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteId := vars["siteId"]

	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.deviceService.AssignDevicesToSite(r.Context(), siteId, req.DeviceIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "assigned", "site": siteId})
}

func (h *AssignmentsHandler) AssignProjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteId := vars["siteId"]

	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.projectService.AssignProjectsToSite(r.Context(), siteId, req.ProjectIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "assigned", "site": siteId})
}

func (h *AssignmentsHandler) AssignDevicesToUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.deviceService.AssignDevicesToUser(r.Context(), userId, req.DeviceIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "assigned", "user": userId})
}
