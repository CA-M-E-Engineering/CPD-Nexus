package services

import (
	"context"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"

	"github.com/google/uuid"
)

type ProjectService struct {
	repo ports.ProjectRepository
}

func NewProjectService(repo ports.ProjectRepository) ports.ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	return s.repo.Get(ctx, id)
}

func (s *ProjectService) ListProjects(ctx context.Context, userID string) ([]domain.Project, error) {
	return s.repo.List(ctx, userID)
}

func (s *ProjectService) CreateProject(ctx context.Context, p *domain.Project) error {
	if p.ID == "" {
		p.ID = "project-" + uuid.New().String()
	}
	return s.repo.Create(ctx, p)
}

func (s *ProjectService) UpdateProject(ctx context.Context, id string, p *domain.Project) error {
	p.ID = id
	return s.repo.Update(ctx, p)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
