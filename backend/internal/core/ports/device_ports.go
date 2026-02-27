package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

// DeviceRepository defines how to store and retrieve devices
type DeviceRepository interface {
	Get(ctx context.Context, userID, id string) (*domain.Device, error)
	GetBySN(ctx context.Context, sn string) (*domain.Device, error)
	List(ctx context.Context, userID string) ([]domain.Device, error)
	ListSNsBySiteID(ctx context.Context, userID, siteID string) ([]string, error)
	Create(ctx context.Context, device *domain.Device) error
	Update(ctx context.Context, device *domain.Device) error
	Delete(ctx context.Context, userID, id string) error

	// Bulk operations
	AssignToUser(ctx context.Context, userID string, deviceIDs []string) error
	AssignToSite(ctx context.Context, siteID string, deviceIDs []string) error
}

// DeviceService defines the business logic for devices
type DeviceService interface {
	GetDevice(ctx context.Context, userID, id string) (*domain.Device, error)
	ListDevices(ctx context.Context, userID string) ([]domain.Device, error)
	RegisterDevice(ctx context.Context, sn, model, userID string) (*domain.Device, error)
	UpdateDevice(ctx context.Context, userID, id string, params map[string]interface{}) error
	DecommissionDevice(ctx context.Context, userID, id string) error

	AssignDevicesToUser(ctx context.Context, userID string, deviceIDs []string) error
	AssignDevicesToSite(ctx context.Context, siteID string, deviceIDs []string) error
}
