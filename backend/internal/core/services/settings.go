package services

import (
	"context"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type SettingsService struct {
	repo ports.SettingsRepository
}

func NewSettingsService(repo ports.SettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) GetSettings(ctx context.Context) (*domain.SystemSettingsResponse, error) {
	settings, err := s.repo.GetSettings(ctx)
	if err != nil {
		return nil, err
	}

	total, deployed, err := s.repo.GetDeviceStats(ctx)
	if err != nil {
		// Log error but maybe don't fail the whole request? For now, let's just return 0s if it fails or return error.
		// Let's return error to be safe.
		return nil, err
	}

	return &domain.SystemSettingsResponse{
		Settings:        *settings,
		TotalDevices:    total,
		DeployedDevices: deployed,
	}, nil
}

func (s *SettingsService) UpdateSettings(ctx context.Context, settings domain.SystemSettings) error {
	// Business logic validation could go here
	return s.repo.UpdateSettings(ctx, settings)
}
