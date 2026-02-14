package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type WorkerRepository interface {
	Get(ctx context.Context, id string) (*domain.Worker, error)
	GetByFIN(ctx context.Context, fin string) (*domain.Worker, error)
	List(ctx context.Context, tenantID, siteID string) ([]domain.Worker, error)
	Create(ctx context.Context, w *domain.Worker) error
	Update(ctx context.Context, w *domain.Worker) error
	Delete(ctx context.Context, id string) error
	GetProjectTenantID(ctx context.Context, projectID string) (string, error)
}

type WorkerService interface {
	GetWorker(ctx context.Context, id string) (*domain.Worker, error)
	ListWorkers(ctx context.Context, tenantID, siteID string) ([]domain.Worker, error)
	CreateWorker(ctx context.Context, w *domain.Worker) error
	UpdateWorker(ctx context.Context, id string, payload map[string]interface{}) error
	DeleteWorker(ctx context.Context, id string) error
}
