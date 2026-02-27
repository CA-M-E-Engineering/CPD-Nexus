package domain

import "time"

type SystemSettings struct {
	ID                   int       `json:"id"`
	AttendanceSyncTime   string    `json:"attendance_sync_time"`    // "HH:MM:SS"
	CPDSubmissionTime    string    `json:"cpd_submission_time"`     // "HH:MM:SS"
	MaxPayloadSizeKB     int       `json:"max_payload_size_kb"`     // KB
	MaxWorkersPerRequest int       `json:"max_workers_per_request"` // Batch size
	MaxRequestsPerMinute int       `json:"max_requests_per_minute"` // Rate limit
	UpdatedAt            time.Time `json:"updated_at"`
}

// DTO to include extra stats not in the settings table
type SystemSettingsResponse struct {
	Settings        SystemSettings `json:"settings"`
	TotalDevices    int            `json:"total_devices"`
	DeployedDevices int            `json:"deployed_devices"`
}
