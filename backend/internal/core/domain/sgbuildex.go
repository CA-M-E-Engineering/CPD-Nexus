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

	// Projects Fields
	ProjectRef          string
	ProjectTitle        string
	ProjectLocation     string
	ProjectContractNo   string
	ProjectContractName string
	HDBPrecinctName     string
	RegulatorID         string
	RegulatorName       string
	OnBehalfOfID        string

	// Main Contractor fields (from Project)
	SiteOwnerName string
	SiteOwnerUEN  string

	// Joined Worker fields
	WorkerName         string
	WorkerFIN          string
	WorkerWorkPassType string
	WorkerNationality  string
	WorkerTrade        string
	EmployerName       string
	EmployerUEN        string
	EmployerTrade      string
	EmployerClientName string
	EmployerClientUEN  string

	// PIC (Person In Charge) for on-behalf submission
	PICName string
	PICFIN  string

	SubmissionDate time.Time
}
