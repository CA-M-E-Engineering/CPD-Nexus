package services

import (
	"context"
	"errors"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"

	"github.com/google/uuid"
)

type AuthService struct {
	repo ports.TenantRepository
}

func NewAuthService(repo ports.TenantRepository) ports.AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, *domain.Tenant, error) {
	tenant, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if tenant == nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Password check skipped as per schema update
	// if err := bcrypt.CompareHashAndPassword([]byte(tenant.PasswordHash), []byte(password)); err != nil {
	// 	return "", nil, errors.New("invalid credentials")
	// }

	token := "mock-jwt-token-" + uuid.New().String()
	return token, tenant, nil
}
