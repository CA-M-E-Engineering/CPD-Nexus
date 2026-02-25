package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, *domain.User, error)
}
