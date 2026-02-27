package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type ProjectRepository interface {
	Get(ctx context.Context, userID, id string) (*domain.Project, error)
	List(ctx context.Context, userID string) ([]domain.Project, error)
	Create(ctx context.Context, p *domain.Project) error
	Update(ctx context.Context, p *domain.Project) error
	Delete(ctx context.Context, userID, id string) error
	AssignToSite(ctx context.Context, siteID string, projectIDs []string) error
}

type ProjectService interface {
	GetProject(ctx context.Context, userID, id string) (*domain.Project, error)
	ListProjects(ctx context.Context, userID string) ([]domain.Project, error)
	CreateProject(ctx context.Context, p *domain.Project) error
	UpdateProject(ctx context.Context, userID, id string, p *domain.Project) error
	DeleteProject(ctx context.Context, userID, id string) error
	AssignProjectsToSite(ctx context.Context, siteID string, projectIDs []string) error
}
