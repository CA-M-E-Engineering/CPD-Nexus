package services

import (
	"context"
	"log"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type SettingsService struct {
	repo            ports.SettingsRepository
	syncScheduler   *DailyScheduler
	submitScheduler *DailyScheduler
}

func NewSettingsService(repo ports.SettingsRepository, sync *DailyScheduler, submit *DailyScheduler) *SettingsService {
	return &SettingsService{
		repo:            repo,
		syncScheduler:   sync,
		submitScheduler: submit,
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
	log.Printf("[SettingsService] Updating system settings in database...")
	if err := s.repo.UpdateSettings(ctx, settings); err != nil {
		return err
	}

	log.Printf("[SettingsService] Database update successful. Notifying schedulers...")

	// Trigger schedulers to re-evaluate their time
	if s.syncScheduler != nil {
		log.Printf("[SettingsService] Resetting AttendanceSync scheduler")
		s.syncScheduler.Reset()
	}
	if s.submitScheduler != nil {
		log.Printf("[SettingsService] Resetting CPDSubmission scheduler")
		s.submitScheduler.Reset()
	}

	return nil
}
