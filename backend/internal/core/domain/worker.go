package domain

type Worker struct {
	ID               string `json:"worker_id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Role             string `json:"role"`
	Status           string `json:"status"`
	TradeCode        string `json:"trade_code,omitempty"`
	CurrentProjectID string `json:"current_project_id,omitempty"`
	CompanyName      string `json:"company_name,omitempty"`
	FIN              string `json:"fin,omitempty"`

	// Details for output (Joined fields)
	ProjectName  string `json:"project_name,omitempty"`
	SiteName     string `json:"site_name,omitempty"`
	SiteLocation string `json:"site_location,omitempty"`
	UserName     string `json:"user_name,omitempty"`
	UserID       string `json:"user_id,omitempty"`
	UserLocation string `json:"user_location,omitempty"`
	UserAddress  string `json:"user_address,omitempty"`
}
