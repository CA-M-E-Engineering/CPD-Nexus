package domain

import (
	"database/sql"
	"time"
)

// AttendanceRow represents a row fetched from the DB with joined info for SGBuildex submission
type AttendanceRow struct {
	AttendanceID string
	DeviceID     string
	WorkerID     string
	SiteID       string
	UserID       string
	TimeIn       time.Time
	TimeOut      sql.NullTime
	Direction    string
	TradeCode    string
	Status       string

	// Joined Site fields
	SiteName     string
	SiteLocation string
	ProjectRef   string

	// Main Contractor fields (from Project)
	SiteOwnerName string
	SiteOwnerUEN  string

	// Joined Worker fields
	WorkerName    string
	WorkerFIN     string
	WorkerTrade   string
	EmployerName  string
	EmployerUEN   string
	EmployerTrade string

	// PIC (Person In Charge) for on-behalf submission
	PICName string
	PICFIN  string

	SubmissionDate time.Time
}
