package services

import (
	"context"
	"fmt"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"

	"github.com/google/uuid"
)

type DeviceService struct {
	repo ports.DeviceRepository
}

func NewDeviceService(repo ports.DeviceRepository) ports.DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) GetDevice(ctx context.Context, id string) (*domain.Device, error) {
	return s.repo.Get(ctx, id)
}

func (s *DeviceService) ListDevices(ctx context.Context, userID string) ([]domain.Device, error) {
	return s.repo.List(ctx, userID)
}

func (s *DeviceService) RegisterDevice(ctx context.Context, sn, model, userID string) (*domain.Device, error) {
	if sn == "" || model == "" {
		return nil, fmt.Errorf("serial number and model are required")
	}

	// Default user if missing (legacy behavior restored per user request)
	if userID == "" {
		userID = "tenant-vendor-1"
	}

	d := &domain.Device{
		ID:     "device-" + uuid.New().String(),
		SN:     sn,
		Model:  model,
		Status: domain.DeviceStatusOffline,
		UserID: userID,
		// Mock site for now, or keep nil
		SiteID: nil, // Let it be null or handled by repository defaults if strict
	}

	// Legacy fallback removed, devices are unassigned to sites by default

	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, id string, params map[string]interface{}) error {
	d, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if d == nil {
		return fmt.Errorf("device not found")
	}

	if v, ok := params["sn"].(string); ok {
		d.SN = v
	}
	if v, ok := params["model"].(string); ok {
		d.Model = v
	}
	if v, ok := params["status"].(string); ok {
		d.Status = domain.DeviceStatus(v)
	}
	if val, ok := params["site_id"]; ok {
		if val == nil {
			d.SiteID = nil
		} else if v, ok := val.(string); ok {
			d.SiteID = &v
		}
	}
	if v, ok := params["user_id"].(string); ok {
		d.UserID = v
	}

	return s.repo.Update(ctx, d)
}

func (s *DeviceService) DecommissionDevice(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *DeviceService) AssignDevicesToUser(ctx context.Context, userID string, deviceIDs []string) error {
	return s.repo.AssignToUser(ctx, userID, deviceIDs)
}
