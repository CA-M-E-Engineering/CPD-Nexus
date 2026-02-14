package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"sgbuildex/internal/core/ports"

	"github.com/gorilla/mux"
)

type DevicesHandler struct {
	Service ports.DeviceService
}

func NewDevicesHandler(service ports.DeviceService) *DevicesHandler {
	return &DevicesHandler{Service: service}
}

func (h *DevicesHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")

	devices, err := h.Service.ListDevices(r.Context(), tenantID)
	if err != nil {
		log.Printf("[GetDevices] Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func (h *DevicesHandler) GetDeviceById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	d, err := h.Service.GetDevice(r.Context(), id)
	if err != nil {
		log.Printf("[GetDeviceById] Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if d == nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *DevicesHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SN       string `json:"sn"`
		Model    string `json:"model"`
		TenantID string `json:"tenant_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := h.Service.RegisterDevice(r.Context(), body.SN, body.Model, body.TenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	if err := h.Service.UpdateDevice(r.Context(), id, params); err != nil {
		log.Printf("[UpdateDevice] Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *DevicesHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.Service.DecommissionDevice(r.Context(), id); err != nil {
		log.Printf("[DeleteDevice] Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "deleted",
		"id":     id,
	})
}
