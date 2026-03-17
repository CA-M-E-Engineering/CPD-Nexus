package services

import (
	"context"
	"fmt"
	"time"

	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/pkg/apperrors"
)

type DeviceService struct {
	repo             ports.DeviceRepository
	analyticsService ports.AnalyticsService
}

func NewDeviceService(repo ports.DeviceRepository, analytics ports.AnalyticsService) ports.DeviceService {
	return &DeviceService{repo: repo, analyticsService: analytics}
}

func (s *DeviceService) GetDevice(ctx context.Context, userID, id string) (*domain.Device, error) {
	if userID == "" {
		return nil, apperrors.NewPermissionDenied("user_id scope required")
	}
	return s.repo.Get(ctx, userID, id)
}

func (s *DeviceService) ListDevices(ctx context.Context, userID, siteID string) ([]domain.Device, error) {
	return s.repo.List(ctx, userID, siteID)
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
		ID:     "d" + time.Now().Format("20060102150405"),
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
	s.analyticsService.LogActivity(ctx, userID, "Device Registered", "device", d.ID, fmt.Sprintf("New device %s (%s) registered", d.SN, d.Model))
	return d, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, userID, id string, params map[string]interface{}) error {
	if userID == "" {
		return apperrors.NewPermissionDenied("user_id scope required")
	}
	d, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return err
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
		if v != d.UserID {
			d.UserID = v
			// Clear site association if owner changes to maintain integrity
			d.SiteID = nil
		}
	}

	err = s.repo.Update(ctx, d)
	if err == nil {
		s.analyticsService.LogActivity(ctx, userID, "Device Updated", "device", id, "Device parameters updated")
	}
	return err
}

func (s *DeviceService) DecommissionDevice(ctx context.Context, userID, id string) error {
	if userID == "" {
		return apperrors.NewPermissionDenied("user_id scope required")
	}
	err := s.repo.Delete(ctx, userID, id)
	if err == nil {
		s.analyticsService.LogActivity(ctx, userID, "Device Decommissioned", "device", id, "Device permanently removed from system")
	}
	return err
}

func (s *DeviceService) AssignDevicesToUser(ctx context.Context, userID string, deviceIDs []string) error {
	if len(deviceIDs) == 0 {
		return nil
	}
	err := s.repo.AssignToUser(ctx, userID, deviceIDs)
	if err == nil {
		actorUserID := ports.GetUserID(ctx)
		// Details: "1 devices" or "5 devices"
		s.analyticsService.LogActivity(ctx, actorUserID, "Device Reassigned", "user", userID, fmt.Sprintf("Assigned %d hardware assets to organization %s", len(deviceIDs), userID))
	}
	return err
}

func (s *DeviceService) AssignDevicesToSite(ctx context.Context, siteID string, deviceIDs []string) error {
	if len(deviceIDs) == 0 {
		return nil
	}
	err := s.repo.AssignToSite(ctx, siteID, deviceIDs)
	if err == nil {
		actorUserID := ports.GetUserID(ctx)
		s.analyticsService.LogActivity(ctx, actorUserID, "Device Reassigned", "site", siteID, fmt.Sprintf("Deployed %d hardware assets to site %s", len(deviceIDs), siteID))
	}
	return err
}
