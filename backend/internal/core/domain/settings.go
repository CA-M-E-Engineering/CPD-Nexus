package domain

import "time"

type SystemSettings struct {
	ID                 int       `json:"id"`
	DeviceSyncInterval string    `json:"device_sync_interval"` // "HH:MM:SS"
	CPDSubmissionTime  string    `json:"cpd_submission_time"`  // "HH:MM:SS"
	ResponseSizeLimit  int64     `json:"response_size_limit"`  // Bytes
	UpdatedAt          time.Time `json:"updated_at"`
}

// DTO to include extra stats not in the settings table
type SystemSettingsResponse struct {
	Settings        SystemSettings `json:"settings"`
	TotalDevices    int            `json:"total_devices"`
	DeployedDevices int            `json:"deployed_devices"`
}
