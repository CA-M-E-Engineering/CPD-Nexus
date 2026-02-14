package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type CompanyRepository interface {
	Get(ctx context.Context, id string) (*domain.Company, error)
	GetByUEN(ctx context.Context, uen string) (*domain.Company, error)
	ListByTenant(ctx context.Context, tenantID string) ([]domain.Company, error)
	Create(ctx context.Context, company *domain.Company) error
	Update(ctx context.Context, company *domain.Company) error
	Delete(ctx context.Context, id string) error
}
