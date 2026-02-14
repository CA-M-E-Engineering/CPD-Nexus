package domain

import "time"

type Company struct {
	ID          string    `json:"company_id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"company_name"`
	UEN         string    `json:"uen"`
	CompanyType string    `json:"company_type"` // contractor | offsite_fabricator
	Address     string    `json:"address,omitempty"`
	Latitude    float64   `json:"latitude,omitempty"`
	Longitude   float64   `json:"longitude,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
