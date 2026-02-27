package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sgbuildex/internal/api/middleware"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/apperrors"

	"github.com/gorilla/mux"
)

type WorkersHandler struct {
	service ports.WorkerService
}

func NewWorkersHandler(service ports.WorkerService) *WorkersHandler {
	return &WorkersHandler{service: service}
}

func (h *WorkersHandler) GetWorkers(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	siteID := r.URL.Query().Get("site_id")

	log.Printf("[WorkersHandler] GetWorkers request: user_id=%s, site_id=%s", userID, siteID)
	workers, err := h.service.ListWorkers(r.Context(), userID, siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workers)
}

func (h *WorkersHandler) GetWorkerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID := middleware.GetUserID(r.Context())
	worker, err := h.service.GetWorker(r.Context(), userID, id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	if worker == nil {
		http.Error(w, "Worker not found", http.StatusNotFound)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(worker)
}

func (h *WorkersHandler) UpdateWorker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r.Context())
	if err := h.service.UpdateWorker(r.Context(), userID, id, payload); err != nil {
		h.handleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *WorkersHandler) DeleteWorker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID := middleware.GetUserID(r.Context())
	if err := h.service.DeleteWorker(r.Context(), userID, id); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
func (h *WorkersHandler) handleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
