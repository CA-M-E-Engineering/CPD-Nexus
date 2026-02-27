package domain

const (
	SyncStatusPendingUpdate       = 0
	SyncStatusSynced              = 1
	SyncStatusPendingRegistration = 2
)

type Worker struct {
	ID               string `json:"worker_id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Role             string `json:"role"`
	UserType         string `json:"user_type"`
	Status           string `json:"status"`
	CurrentProjectID string `json:"current_project_id,omitempty"`

	// SGBuildex Compliance Fields
	PersonIDNo              string `json:"person_id_no,omitempty"`
	PersonIDAndWorkPassType string `json:"person_id_and_work_pass_type,omitempty"`
	PersonNationality       string `json:"person_nationality,omitempty"`
	PersonTrade             string `json:"person_trade,omitempty"`

	// IoT Authentication Fields
	AuthStartTime string `json:"auth_start_time,omitempty"`
	AuthEndTime   string `json:"auth_end_time,omitempty"`
	FDID          int    `json:"fdid,omitempty"`
	FaceImgLoc    string `json:"face_img_loc,omitempty"`
	CardNumber    string `json:"card_number,omitempty"`
	CardType      string `json:"card_type,omitempty"`
	IsSynced      int    `json:"is_synced"`

	// Details for output (Joined fields)
	ProjectName  string `json:"project_name,omitempty"`
	SiteID       string `json:"site_id,omitempty"`
	SiteName     string `json:"site_name,omitempty"`
	SiteLocation string `json:"site_location,omitempty"`
	UserName     string `json:"user_name,omitempty"`
	UserID       string `json:"user_id,omitempty"`
	UserLocation string `json:"user_location,omitempty"`
	UserAddress  string `json:"user_address,omitempty"`
}
