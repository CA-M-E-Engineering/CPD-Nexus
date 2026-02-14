package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type SettingsService interface {
	GetSettings(ctx context.Context) (*domain.SystemSettingsResponse, error)
	UpdateSettings(ctx context.Context, settings domain.SystemSettings) error
}

type SettingsRepository interface {
	GetSettings(ctx context.Context) (*domain.SystemSettings, error)
	UpdateSettings(ctx context.Context, settings domain.SystemSettings) error
	GetDeviceStats(ctx context.Context) (total int, online int, err error)
}
