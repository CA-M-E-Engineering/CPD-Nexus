package domain

import (
	"time"
)

// DeviceStatus represents the operational status of a device
type DeviceStatus string

const (
	DeviceStatusOnline   DeviceStatus = "online"
	DeviceStatusOffline  DeviceStatus = "offline"
	DeviceStatusInactive DeviceStatus = "inactive"
	DeviceStatusUnknown  DeviceStatus = "unknown"
)

// Device represents an IoT device in the system
type Device struct {
	ID     string       `json:"device_id"`
	SN     string       `json:"sn"`
	Model  string       `json:"model"`
	Status DeviceStatus `json:"status"`

	// Associations
	SiteID   *string `json:"site_id,omitempty"` // Pointer to allow NULL
	TenantID string  `json:"tenant_id"`         // Required

	// Metadata (often joined)
	SiteName   string `json:"site_name,omitempty"`
	TenantName string `json:"tenant_name,omitempty"`

	// Telemetry
	LastHeartbeat   *time.Time `json:"last_heartbeat,omitempty"`
	LastOnlineCheck *time.Time `json:"last_online_check,omitempty"`
	Battery         int        `json:"battery"`
}
