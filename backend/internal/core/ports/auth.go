package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, *domain.User, error)
}
