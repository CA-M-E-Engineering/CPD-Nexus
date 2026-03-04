package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

// PitstopRepository defines the data access methods for Pitstop configurations
type PitstopRepository interface {
	GetAuthorisations(ctx context.Context) ([]*domain.PitstopAuthorisation, error)
	InsertAuthorisations(ctx context.Context, auths []*domain.PitstopAuthorisation) error
	UpdateAuthorisations(ctx context.Context, auths []*domain.PitstopAuthorisation) error
	AssignOnBehalfOfToUser(ctx context.Context, userID string, onBehalfOfNames []string) error
}

// PitstopService defines the use-case operations for Pitstop/SGBuildex integration
type PitstopService interface {
	GetAuthorisations(ctx context.Context) ([]*domain.PitstopAuthorisation, error)
	SyncConfig(ctx context.Context, userID string) error
	TestSubmission(ctx context.Context, projectID string) (int, error)
	GetProjectsWithPendingAttendance(ctx context.Context) ([]domain.Project, error)
	SubmitPendingAttendance(ctx context.Context) error
	AssignOnBehalfOfToUser(ctx context.Context, userID string, onBehalfOfNames []string) error
}
