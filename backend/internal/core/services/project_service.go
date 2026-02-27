package services

import (
	"context"
	"time"

	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/apperrors"
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
	if p.ID == "" {
		p.ID = "p" + time.Now().Format("20060102150405")
	}
	return s.repo.Create(ctx, p)
}

func (s *ProjectService) UpdateProject(ctx context.Context, userID, id string, p *domain.Project) error {
	if userID == "" {
		return apperrors.NewPermissionDenied("user_id scope required")
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
