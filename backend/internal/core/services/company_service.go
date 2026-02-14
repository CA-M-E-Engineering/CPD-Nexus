package services

import (
	"context"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type CompanyService struct {
	repo ports.CompanyRepository
}

func NewCompanyService(repo ports.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) GetCompany(ctx context.Context, id string) (*domain.Company, error) {
	return s.repo.Get(ctx, id)
}

func (s *CompanyService) GetCompanyByUEN(ctx context.Context, uen string) (*domain.Company, error) {
	return s.repo.GetByUEN(ctx, uen)
}

func (s *CompanyService) ListCompaniesByTenant(ctx context.Context, tenantID string) ([]domain.Company, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}

func (s *CompanyService) CreateCompany(ctx context.Context, company *domain.Company) error {
	return s.repo.Create(ctx, company)
}

func (s *CompanyService) UpdateCompany(ctx context.Context, company *domain.Company) error {
	return s.repo.Update(ctx, company)
}

func (s *CompanyService) DeleteCompany(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
