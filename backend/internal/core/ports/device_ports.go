package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

// DeviceRepository defines how to store and retrieve devices
type DeviceRepository interface {
	Get(ctx context.Context, id string) (*domain.Device, error)
	GetBySN(ctx context.Context, sn string) (*domain.Device, error)
	List(ctx context.Context, tenantID string) ([]domain.Device, error)
	Create(ctx context.Context, device *domain.Device) error
	Update(ctx context.Context, device *domain.Device) error
	Delete(ctx context.Context, id string) error

	// Bulk operations
	AssignToTenant(ctx context.Context, tenantID string, deviceIDs []string) error
}

// DeviceService defines the business logic for devices
type DeviceService interface {
	GetDevice(ctx context.Context, id string) (*domain.Device, error)
	ListDevices(ctx context.Context, tenantID string) ([]domain.Device, error)
	RegisterDevice(ctx context.Context, sn, model, tenantID string) (*domain.Device, error)
	UpdateDevice(ctx context.Context, id string, params map[string]interface{}) error
	DecommissionDevice(ctx context.Context, id string) error

	AssignDevicesToTenant(ctx context.Context, tenantID string, deviceIDs []string) error
}
