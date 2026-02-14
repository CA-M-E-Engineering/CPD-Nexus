package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type TenantRepository interface {
	GetByUsername(ctx context.Context, username string) (*domain.Tenant, error)
	Get(ctx context.Context, id string) (*domain.Tenant, error)
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, *domain.Tenant, error)
}
