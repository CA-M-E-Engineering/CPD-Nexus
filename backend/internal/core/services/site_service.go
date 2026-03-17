package services

import (
	"context"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
)

type SiteService struct {
	repo      ports.SiteRepository
	analytics ports.AnalyticsService
}

func NewSiteService(repo ports.SiteRepository, analytics ports.AnalyticsService) ports.SiteService {
	return &SiteService{
		repo:      repo,
		analytics: analytics,
	}
}

func (s *SiteService) GetSite(ctx context.Context, userID, id string) (*domain.Site, error) {
	return s.repo.Get(ctx, userID, id)
}

func (s *SiteService) ListSites(ctx context.Context, userID string) ([]domain.Site, error) {
	return s.repo.List(ctx, userID)
}

func (s *SiteService) CreateSite(ctx context.Context, site *domain.Site) error {
	err := s.repo.Create(ctx, site)
	if err == nil {
		s.analytics.LogActivity(ctx, site.UserID, "Site Created", "site", site.ID, "New site "+site.Name+" established")
	}
	return err
}

func (s *SiteService) UpdateSite(ctx context.Context, userID, id string, site *domain.Site) error {
	err := s.repo.Update(ctx, site)
	if err == nil {
		s.analytics.LogActivity(ctx, userID, "Site Updated", "site", id, "Site details for "+site.Name+" modified")
	}
	return err
}

func (s *SiteService) DeleteSite(ctx context.Context, userID, id string) error {
	err := s.repo.Delete(ctx, userID, id)
	if err == nil {
		s.analytics.LogActivity(ctx, userID, "Site Deleted", "site", id, "Site record removed from system")
	}
	return err
}
