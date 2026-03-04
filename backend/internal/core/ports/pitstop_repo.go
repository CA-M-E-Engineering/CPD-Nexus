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
