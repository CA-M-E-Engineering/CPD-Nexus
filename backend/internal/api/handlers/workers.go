package handlers

import (
	"encoding/json"
	"net/http"

	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"

	"github.com/gorilla/mux"
)

type WorkersHandler struct {
	service ports.WorkerService
}

func NewWorkersHandler(service ports.WorkerService) *WorkersHandler {
	return &WorkersHandler{service: service}
}

func (h *WorkersHandler) GetWorkers(w http.ResponseWriter, r *http.Request) {
	// userID MUST come from the middleware context, not the query string,
	// to enforce multi-tenant isolation.
	userID := ports.GetUserID(r.Context())
	siteID := r.URL.Query().Get("site_id")

	workers, err := h.service.ListWorkers(r.Context(), userID, siteID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workers)
}

func (h *WorkersHandler) GetWorkerById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID := ports.GetUserID(r.Context())

	worker, err := h.service.GetWorker(r.Context(), userID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	if worker == nil {
		http.Error(w, "worker not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(worker)
}

func (h *WorkersHandler) CreateWorker(w http.ResponseWriter, r *http.Request) {
	var worker domain.Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateWorker(r.Context(), &worker); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(worker)
}

func (h *WorkersHandler) UpdateWorker(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID := ports.GetUserID(r.Context())

	var req domain.UpdateWorkerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateWorker(r.Context(), userID, id, &req); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *WorkersHandler) DeleteWorker(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID := ports.GetUserID(r.Context())

	if err := h.service.DeleteWorker(r.Context(), userID, id); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
