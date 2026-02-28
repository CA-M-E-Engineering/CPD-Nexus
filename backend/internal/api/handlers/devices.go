package handlers

import (
	"encoding/json"
	"net/http"

	"sgbuildex/internal/api/middleware"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/apperrors"

	"github.com/gorilla/mux"
)

type DevicesHandler struct {
	Service ports.DeviceService
}

func NewDevicesHandler(service ports.DeviceService) *DevicesHandler {
	return &DevicesHandler{Service: service}
}

func (h *DevicesHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	devices, err := h.Service.ListDevices(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func (h *DevicesHandler) GetDeviceById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID := middleware.GetUserID(r.Context())
	d, err := h.Service.GetDevice(r.Context(), userID, id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *DevicesHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SN     string `json:"sn"`
		Model  string `json:"model"`
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := h.Service.RegisterDevice(r.Context(), body.SN, body.Model, body.UserID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

func (h *DevicesHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r.Context())
	if err := h.Service.UpdateDevice(r.Context(), userID, id, params); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *DevicesHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID := middleware.GetUserID(r.Context())
	if err := h.Service.DecommissionDevice(r.Context(), userID, id); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "deleted",
		"id":     id,
	})
}
func (h *DevicesHandler) handleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
