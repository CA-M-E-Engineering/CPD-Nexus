package services

import (
	"context"
	"time"

	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/pkg/apperrors"
	"cpd-nexus/internal/pkg/validation"
)

type ProjectService struct {
	repo ports.ProjectRepository
}

func NewProjectService(repo ports.ProjectRepository) ports.ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetProject(ctx context.Context, userID, id string) (*domain.Project, error) {
	if userID == "" {
		return nil, apperrors.NewPermissionDenied("user_id scope required")
	}
	return s.repo.Get(ctx, userID, id)
}

func (s *ProjectService) ListProjects(ctx context.Context, userID string) ([]domain.Project, error) {
	return s.repo.List(ctx, userID)
}

func (s *ProjectService) CreateProject(ctx context.Context, p *domain.Project) error {
	if err := s.validateProject(p); err != nil {
		return err
	}
	if p.ID == "" {
		p.ID = "p" + time.Now().Format("20060102150405")
	}
	return s.repo.Create(ctx, p)
}

func (s *ProjectService) validateProject(p *domain.Project) error {
	if len(p.Title) > 1000 {
		return apperrors.NewValidationError("project_title too long (max 1000)")
	}
	if p.Reference != "" && !validation.ValidateProjectReferenceNumber(p.Reference) {
		return apperrors.NewValidationError("invalid project_reference_number format (e.g. A1234-12345-2022)")
	}
	if len(p.Location) > 2000 {
		return apperrors.NewValidationError("project_location too long (max 2000)")
	}
	if p.ContractRef != "" {
		// Try to validate as either HDB or LTA. If neither matches, might be okay but warn?
		// User rules say HDB or LTA mandatory if participant is active.
		if !validation.ValidateHDBContractNumber(p.ContractRef) && !validation.ValidateLTAContractNumber(p.ContractRef) {
			// Note: We are strict here as per user request to "enforce"
			return apperrors.NewValidationError("invalid project_contract_number format (HDB: D/NNNNN/YY, LTA: Max 20 chars)")
		}
	}
	if len(p.ContractName) > 100 {
		return apperrors.NewValidationError("project_contract_name too long (max 100)")
	}
	if len(p.HDBPrecinct) > 40 {
		return apperrors.NewValidationError("hdb_precinct_name too long (max 40)")
	}
	if p.MainContractorUEN != "" && !validation.ValidateUEN(p.MainContractorUEN) {
		return apperrors.NewValidationError("invalid main_contractor_uen format")
	}
	if p.WorkerCompanyUEN != "" && !validation.ValidateUEN(p.WorkerCompanyUEN) {
		return apperrors.NewValidationError("invalid worker_company_uen format")
	}
	if p.WorkerCompanyClientUEN != "" && !validation.ValidateUEN(p.WorkerCompanyClientUEN) {
		return apperrors.NewValidationError("invalid worker_company_client_uen format")
	}
	return nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, userID, id string, p *domain.Project) error {
	if userID == "" {
		return apperrors.NewPermissionDenied("user_id scope required")
	}
	if err := s.validateProject(p); err != nil {
		return err
	}
	// Verify ownership before update
	existing, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return err
	}
	p.ID = existing.ID
	p.UserID = existing.UserID
	return s.repo.Update(ctx, p)
}

func (s *ProjectService) DeleteProject(ctx context.Context, userID, id string) error {
	if userID == "" {
		return apperrors.NewPermissionDenied("user_id scope required")
	}
	return s.repo.Delete(ctx, userID, id)
}

func (s *ProjectService) AssignProjectsToSite(ctx context.Context, siteID string, projectIDs []string) error {
	return s.repo.AssignToSite(ctx, siteID, projectIDs)
}
