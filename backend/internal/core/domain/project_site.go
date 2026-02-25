package domain

import "time"

type Project struct {
	ID           string `json:"project_id"`
	SiteID       string `json:"site_id"`
	UserID       string `json:"user_id"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	Reference    string `json:"reference"`
	ContractRef  string `json:"contract"`
	ContractName string `json:"contract_name"`
	Location     string `json:"location"`
	HDBPrecinct  string `json:"hdb_precinct"`

	// Inline company details (no FK to companies table)
	MainContractorName        string `json:"main_contractor_name,omitempty"`
	MainContractorUEN         string `json:"main_contractor_uen,omitempty"`
	OffsiteFabricatorName     string `json:"offsite_fabricator_name,omitempty"`
	OffsiteFabricatorUEN      string `json:"offsite_fabricator_uen,omitempty"`
	OffsiteFabricatorLocation string `json:"offsite_fabricator_location,omitempty"`
	WorkerCompanyName         string `json:"worker_company_name,omitempty"`
	WorkerCompanyUEN          string `json:"worker_company_uen,omitempty"`
	WorkerCompanyClientName   string `json:"worker_company_client_name,omitempty"`
	WorkerCompanyClientUEN    string `json:"worker_company_client_uen,omitempty"`
	WorkerCompanyTrade        string `json:"worker_company_trade,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Calculated fields
	WorkerCount int    `json:"worker_count"`
	DeviceCount int    `json:"device_count"`
	SiteName    string `json:"site_name,omitempty"`
}

type Site struct {
	ID        string    `json:"site_id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"site_name"`
	Location  string    `json:"location,omitempty"`
	Latitude  float64   `json:"lat"`
	Longitude float64   `json:"lng"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Calculated fields
	WorkerCount int    `json:"worker_count"`
	DeviceCount int    `json:"device_count"`
	UserName    string `json:"user_name,omitempty"`
}
