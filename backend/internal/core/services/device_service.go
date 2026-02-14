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

func (s *DeviceService) ListDevices(ctx context.Context, tenantID string) ([]domain.Device, error) {
	return s.repo.List(ctx, tenantID)
}

func (s *DeviceService) RegisterDevice(ctx context.Context, sn, model, tenantID string) (*domain.Device, error) {
	if sn == "" || model == "" {
		return nil, fmt.Errorf("serial number and model are required")
	}

	// Default tenant if missing (legacy behavior)
	if tenantID == "" {
		tenantID = "tenant-uuid-1"
	}

	d := &domain.Device{
		ID:       "device-" + uuid.New().String(),
		SN:       sn,
		Model:    model,
		Status:   domain.DeviceStatusOffline,
		TenantID: tenantID,
		// Mock site for now, or keep nil
		SiteID: nil, // Let it be null or handled by repository defaults if strict
	}

	// Legacy fallback: if DB constraint requires site_id being 'site-uuid-1'
	// For now let's assume NULL is allowed or we set it
	defaultSite := "site-uuid-1"
	d.SiteID = &defaultSite

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
	if v, ok := params["site_id"].(string); ok {
		d.SiteID = &v
	}
	if v, ok := params["tenant_id"].(string); ok {
		d.TenantID = v
	}

	return s.repo.Update(ctx, d)
}

func (s *DeviceService) DecommissionDevice(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *DeviceService) AssignDevicesToTenant(ctx context.Context, tenantID string, deviceIDs []string) error {
	return s.repo.AssignToTenant(ctx, tenantID, deviceIDs)
}
