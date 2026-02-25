package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type UserRepository interface {
	Get(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type UserService interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	CreateUser(ctx context.Context, user *domain.User, password string) error
	UpdateUser(ctx context.Context, id string, payload map[string]interface{}) error
	DeleteUser(ctx context.Context, id string) error
}
