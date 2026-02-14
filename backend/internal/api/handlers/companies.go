package handlers

import (
	"encoding/json"
	"net/http"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/services"
	"sgbuildex/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CompaniesHandler struct {
	service *services.CompanyService
}

func NewCompaniesHandler(service *services.CompanyService) *CompaniesHandler {
	return &CompaniesHandler{service: service}
}

func (h *CompaniesHandler) GetCompanyById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	company, err := h.service.GetCompany(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if company == nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

func (h *CompaniesHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "tenant_id query parameter is required", http.StatusBadRequest)
		return
	}

	companies, err := h.service.ListCompaniesByTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

func (h *CompaniesHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// DEBUG LOG
	logger.Infof("--- CREATE COMPANY DB CHECK ---")
	logger.Infof("Payload TenantID: [%s]", company.TenantID)
	logger.Infof("-------------------------------")

	if company.ID == "" {
		company.ID = "company-" + uuid.New().String()
	}

	if company.Status == "" {
		company.Status = "active"
	}

	if err := h.service.CreateCompany(r.Context(), &company); err != nil {
		logger.Errorf("CREATE COMPANY FAILED: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(company)
}

func (h *CompaniesHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	company.ID = id

	if err := h.service.UpdateCompany(r.Context(), &company); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

func (h *CompaniesHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteCompany(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
