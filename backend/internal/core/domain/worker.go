package domain

type Worker struct {
	ID               string `json:"worker_id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Role             string `json:"role"`
	Status           string `json:"status"`
	CurrentProjectID string `json:"current_project_id,omitempty"`

	// SGBuildex Compliance Fields
	PersonIDNo              string `json:"person_id_no,omitempty"`
	PersonIDAndWorkPassType string `json:"person_id_and_work_pass_type,omitempty"`
	PersonNationality       string `json:"person_nationality,omitempty"`
	PersonTrade             string `json:"person_trade,omitempty"`

	// Details for output (Joined fields)
	ProjectName  string `json:"project_name,omitempty"`
	SiteName     string `json:"site_name,omitempty"`
	SiteLocation string `json:"site_location,omitempty"`
	UserName     string `json:"user_name,omitempty"`
	UserID       string `json:"user_id,omitempty"`
	UserLocation string `json:"user_location,omitempty"`
	UserAddress  string `json:"user_address,omitempty"`
}
