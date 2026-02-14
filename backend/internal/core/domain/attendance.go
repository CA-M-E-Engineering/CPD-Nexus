package domain

import "time"

type Attendance struct {
	ID              string     `json:"attendance_id"`
	DeviceID        string     `json:"device_id"`
	WorkerID        string     `json:"worker_id"`
	SiteID          string     `json:"site_id"`
	TenantID        string     `json:"tenant_id"`
	TimeIn          *time.Time `json:"time_in"`
	TimeOut         *time.Time `json:"time_out"`
	Direction       string     `json:"direction"`
	TradeCode       string     `json:"trade_code"`
	Status          string     `json:"status"`
	SubmissionDate  string     `json:"submission_date"`
	ResponsePayload string     `json:"response_payload,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Joined fields
	WorkerName string `json:"worker_name,omitempty"`
	SiteName   string `json:"site_name,omitempty"`
}
