package ports

import (
	"context"
	"cpd-nexus/internal/core/domain"
)

// PitstopRepository defines the data access methods for Pitstop configurations
type PitstopRepository interface {
	GetAuthorisations(ctx context.Context, userID string) ([]*domain.PitstopAuthorisation, error)
	InsertAuthorisations(ctx context.Context, auths []*domain.PitstopAuthorisation) error
	UpdateAuthorisations(ctx context.Context, auths []*domain.PitstopAuthorisation) error
	AssignOnBehalfOfToUser(ctx context.Context, userID string, onBehalfOfNames []string) error
}

// PitstopService defines the use-case operations for Pitstop/SGBuildex integration
type PitstopService interface {
	GetAuthorisations(ctx context.Context, userID string) ([]*domain.PitstopAuthorisation, error)
	SyncConfig(ctx context.Context, userID string) error
	TestSubmission(ctx context.Context, userID, projectID string) (int, int, error)
	GetProjectsWithPendingAttendance(ctx context.Context, userID string) ([]domain.Project, error)
	SubmitPendingAttendance(ctx context.Context) error
	AssignOnBehalfOfToUser(ctx context.Context, userID string, onBehalfOfNames []string) error
}
