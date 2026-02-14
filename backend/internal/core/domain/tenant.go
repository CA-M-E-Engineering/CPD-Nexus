package domain

type Tenant struct {
	ID         string `json:"tenant_id"`
	Name       string `json:"tenant_name"`
	Username   string `json:"username"`
	TenantType string `json:"tenant_type"`
	Status     string `json:"status"`
	Latitude   string `json:"latitude,omitempty"`
	Longitude  string `json:"longitude,omitempty"`
}
