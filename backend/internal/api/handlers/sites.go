package handlers

import (
	"encoding/json"
	"net/http"

	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"

	"github.com/gorilla/mux"
)

type SitesHandler struct {
	service ports.SiteService
}

func NewSitesHandler(service ports.SiteService) *SitesHandler {
	return &SitesHandler{service: service}
}

func (h *SitesHandler) GetSites(w http.ResponseWriter, r *http.Request) {
	userID := ports.GetUserID(r.Context())
	sites, err := h.service.ListSites(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

func (h *SitesHandler) GetSiteById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID := ports.GetUserID(r.Context())

	site, err := h.service.GetSite(r.Context(), userID, id)
	if err != nil {
		writeError(w, err)
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

	// Enforce multi-tenancy: set UserID from authenticated context
	site.UserID = ports.GetUserID(r.Context())
	
	if err := h.service.CreateSite(r.Context(), &site); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}

func (h *SitesHandler) UpdateSite(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID := ports.GetUserID(r.Context())

	var site domain.Site
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateSite(r.Context(), userID, id, &site); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *SitesHandler) DeleteSite(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID := ports.GetUserID(r.Context())

	if err := h.service.DeleteSite(r.Context(), userID, id); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
