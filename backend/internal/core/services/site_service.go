package services

import (
	"context"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"

	"github.com/google/uuid"
)

type SiteService struct {
	repo ports.SiteRepository
}

func NewSiteService(repo ports.SiteRepository) ports.SiteService {
	return &SiteService{repo: repo}
}

func (s *SiteService) GetSite(ctx context.Context, id string) (*domain.Site, error) {
	return s.repo.Get(ctx, id)
}

func (s *SiteService) ListSites(ctx context.Context, userID string) ([]domain.Site, error) {
	return s.repo.List(ctx, userID)
}

func (s *SiteService) CreateSite(ctx context.Context, site *domain.Site) error {
	if site.ID == "" {
		site.ID = "site-" + uuid.New().String()
	}
	return s.repo.Create(ctx, site)
}

func (s *SiteService) UpdateSite(ctx context.Context, id string, site *domain.Site) error {
	site.ID = id
	return s.repo.Update(ctx, site)
}

func (s *SiteService) DeleteSite(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
