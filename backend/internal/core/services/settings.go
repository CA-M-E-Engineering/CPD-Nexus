package services

import (
	"context"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/pkg/logger"
)

type SettingsService struct {
	repo            ports.SettingsRepository
	syncScheduler   *DailyScheduler
	submitScheduler *DailyScheduler
	analytics       ports.AnalyticsService
}

func NewSettingsService(repo ports.SettingsRepository, sync *DailyScheduler, submit *DailyScheduler, analytics ports.AnalyticsService) ports.SettingsService {
	return &SettingsService{
		repo:            repo,
		syncScheduler:   sync,
		submitScheduler: submit,
		analytics:       analytics,
	}
}

func (s *SettingsService) GetSettings(ctx context.Context) (*domain.SystemSettingsResponse, error) {
	settings, err := s.repo.GetSettings(ctx)
	if err != nil {
		return nil, err
	}

	total, deployed, err := s.repo.GetDeviceStats(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.SystemSettingsResponse{
		Settings:        *settings,
		TotalDevices:    total,
		DeployedDevices: deployed,
	}, nil
}

func (s *SettingsService) UpdateSettings(ctx context.Context, settings domain.SystemSettings) error {
	logger.Infof("[SettingsService] Updating system settings in database...")
	if err := s.repo.UpdateSettings(ctx, settings); err != nil {
		return err
	}

	logger.Infof("[SettingsService] Database update successful. Notifying schedulers...")

	// Log activity
	userID := ports.GetUserID(ctx)
	s.analytics.LogActivity(ctx, userID, "Settings Updated", "system", "global", "System-wide parameters and schedules modified")

	// Trigger schedulers to re-evaluate their time
	if s.syncScheduler != nil {
		logger.Infof("[SettingsService] Resetting AttendanceSync scheduler")
		s.syncScheduler.Reset()
	}
	if s.submitScheduler != nil {
		logger.Infof("[SettingsService] Resetting CPDSubmission scheduler")
		s.submitScheduler.Reset()
	}

	return nil
}
