package services

import (
	"context"
	"time"

	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/apperrors"
)

type SiteService struct {
	repo ports.SiteRepository
}

func NewSiteService(repo ports.SiteRepository) ports.SiteService {
	return &SiteService{repo: repo}
}

func (s *SiteService) GetSite(ctx context.Context, userID, id string) (*domain.Site, error) {
	if userID == "" {
		return nil, apperrors.NewPermissionDenied("user_id scope required")
	}
	return s.repo.Get(ctx, userID, id)
}

func (s *SiteService) ListSites(ctx context.Context, userID string) ([]domain.Site, error) {
	return s.repo.List(ctx, userID)
}

func (s *SiteService) CreateSite(ctx context.Context, site *domain.Site) error {
	if site.ID == "" {
		site.ID = "s" + time.Now().Format("20060102150405")
	}
	return s.repo.Create(ctx, site)
}

func (s *SiteService) UpdateSite(ctx context.Context, userID, id string, site *domain.Site) error {
	if userID == "" {
		return apperrors.NewPermissionDenied("user_id scope required")
	}
	// Verify ownership before update
	existing, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return err
	}
	site.ID = existing.ID
	site.UserID = existing.UserID
	return s.repo.Update(ctx, site)
}

func (s *SiteService) DeleteSite(ctx context.Context, userID, id string) error {
	if userID == "" {
		return apperrors.NewPermissionDenied("user_id scope required")
	}
	return s.repo.Delete(ctx, userID, id)
}
