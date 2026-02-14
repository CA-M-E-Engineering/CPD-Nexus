package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type SiteRepository interface {
	Get(ctx context.Context, id string) (*domain.Site, error)
	List(ctx context.Context, userID string) ([]domain.Site, error)
	Create(ctx context.Context, s *domain.Site) error
	Update(ctx context.Context, s *domain.Site) error
	Delete(ctx context.Context, id string) error
}

type SiteService interface {
	GetSite(ctx context.Context, id string) (*domain.Site, error)
	ListSites(ctx context.Context, userID string) ([]domain.Site, error)
	CreateSite(ctx context.Context, s *domain.Site) error
	UpdateSite(ctx context.Context, id string, s *domain.Site) error
	DeleteSite(ctx context.Context, id string) error
}
