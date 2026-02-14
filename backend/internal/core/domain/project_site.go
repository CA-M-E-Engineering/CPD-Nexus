package domain

import "time"

type Project struct {
	ID                    string    `json:"project_id"`
	SiteID                string    `json:"site_id"`
	TenantID              string    `json:"tenant_id"`
	Title                 string    `json:"title"`
	Status                string    `json:"status"`
	Reference             string    `json:"reference"`
	ContractRef           string    `json:"contract"`
	ContractName          string    `json:"contract_name"`
	Location              string    `json:"location"`
	HDBPrecinct           string    `json:"hdb_precinct"`
	MainContractorID      string    `json:"main_contractor_id,omitempty"`
	OffsiteFabricatorID   string    `json:"offsite_fabricator_id,omitempty"`
	WorkerCompanyID       string    `json:"worker_company_id,omitempty"`
	WorkerCompanyClientID string    `json:"worker_company_client_id,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	// Company Names (Resolved)
	MainContractorName      string `json:"main_contractor_name,omitempty"`
	OffsiteFabricatorName   string `json:"offsite_fabricator_name,omitempty"`
	WorkerCompanyName       string `json:"worker_company_name,omitempty"`
	WorkerCompanyClientName string `json:"worker_company_client_name,omitempty"`

	// Calculated fields
	WorkerCount int    `json:"worker_count"`
	DeviceCount int    `json:"device_count"`
	SiteName    string `json:"site_name,omitempty"`
}

type Site struct {
	ID        string    `json:"site_id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"site_name"`
	Location  string    `json:"location,omitempty"`
	Latitude  float64   `json:"lat"`
	Longitude float64   `json:"lng"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Calculated fields
	WorkerCount int    `json:"worker_count"`
	DeviceCount int    `json:"device_count"`
	TenantName  string `json:"tenant_name,omitempty"`
}
