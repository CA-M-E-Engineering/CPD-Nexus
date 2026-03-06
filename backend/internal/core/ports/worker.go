package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

// IsVendorContextKey is used by the service layer to pass vendor/admin status
// into the repository layer via context — avoiding a dependency on the HTTP middleware package.
type workerContextKey string

const IsVendorContextKey workerContextKey = "workerIsVendor"

type WorkerRepository interface {
	Get(ctx context.Context, userID, id string) (*domain.Worker, error)
	GetByFIN(ctx context.Context, fin string) (*domain.Worker, error)
	List(ctx context.Context, userID, siteID string) ([]domain.Worker, error)
	ListByIsSynced(ctx context.Context, userID string, syncStatus int) ([]domain.Worker, error)
	Create(ctx context.Context, w *domain.Worker) error
	Update(ctx context.Context, w *domain.Worker) error
	MarkSynced(ctx context.Context, id string) error
	Delete(ctx context.Context, userID, id string) error
	GetProjectUserID(ctx context.Context, projectID string) (string, error)
	AssignToProject(ctx context.Context, projectID string, workerIDs []string, userID string) error
}

type WorkerService interface {
	GetWorker(ctx context.Context, userID, id string) (*domain.Worker, error)
	ListWorkers(ctx context.Context, userID, siteID string) ([]domain.Worker, error)
	ListPendingSyncWorkers(ctx context.Context, userID string) ([]domain.Worker, error)
	CreateWorker(ctx context.Context, w *domain.Worker) error
	UpdateWorker(ctx context.Context, userID, id string, payload map[string]interface{}) error
	DeleteWorker(ctx context.Context, userID, id string) error
	AssignWorkersToProject(ctx context.Context, projectID string, workerIDs []string) error
}
