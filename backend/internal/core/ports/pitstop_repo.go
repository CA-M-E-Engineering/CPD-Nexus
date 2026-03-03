package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

// PitstopRepository defines the data access methods for Pitstop configurations
type PitstopRepository interface {
	GetAuthorisations(ctx context.Context) ([]*domain.PitstopAuthorisation, error)
	UpsertAuthorisations(ctx context.Context, auths []*domain.PitstopAuthorisation) error
}
