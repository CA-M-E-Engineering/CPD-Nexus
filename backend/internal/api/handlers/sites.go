package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"

	"github.com/gorilla/mux"
)

type SitesHandler struct {
	service ports.SiteService
}

func NewSitesHandler(service ports.SiteService) *SitesHandler {
	return &SitesHandler{service: service}
}

func (h *SitesHandler) GetSites(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	sites, err := h.service.ListSites(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

func (h *SitesHandler) GetSiteById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	site, err := h.service.GetSite(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if site == nil {
		http.Error(w, "Site not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}

func (h *SitesHandler) CreateSite(w http.ResponseWriter, r *http.Request) {
	var site domain.Site
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateSite(r.Context(), &site); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}

func (h *SitesHandler) UpdateSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var site domain.Site
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateSite(r.Context(), id, &site); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *SitesHandler) DeleteSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteSite(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
