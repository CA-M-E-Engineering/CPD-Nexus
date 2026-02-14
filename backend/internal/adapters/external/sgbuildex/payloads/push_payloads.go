package payloads

// ManpowerUtilization represents the manpower utilization record for a project
type ManpowerUtilization struct {
	// Internal fields (not exported to JSON)
	InternalAttendanceID string `json:"-"`
	InternalWorkerID     string `json:"-"`
	InternalSiteID       string `json:"-"`

	SubmissionEntity int    `json:"submission_entity"`
	SubmissionMonth  string `json:"submission_month"` // YYYY-MM

	// Onsite Builder (submission_entity = 1)
	ProjectReferenceNumber     *string `json:"project_reference_number,omitempty"`
	ProjectTitle               *string `json:"project_title,omitempty"`
	ProjectLocationDescription *string `json:"project_location_description,omitempty"`
	MainContractorCompanyName  *string `json:"main_contractor_company_name,omitempty"`
	MainContractorCompanyUEN   *string `json:"main_contractor_company_unique_entity_number,omitempty"`

	// Offsite Fabricator (submission_entity = 2)
	OffsiteFabricatorCompanyName         *string `json:"offsite_fabricator_company_name,omitempty"`
	OffsiteFabricatorCompanyUEN          *string `json:"offsite_fabricator_company_unique_entity_number,omitempty"`
	OffsiteFabricatorLocationDescription *string `json:"offsite_fabricator_location_description,omitempty"`

	// Person
	PersonIDNo                      string   `json:"person_id_no"`
	PersonName                      string   `json:"person_name"`
	PersonIDAndWorkPassType         string   `json:"person_id_and_work_pass_type"`
	PersonTrade                     string   `json:"person_trade"`
	PersonEmployerCompanyName       string   `json:"person_employer_company_name"`
	PersonEmployerCompanyUEN        string   `json:"person_employer_company_unique_entity_number"`
	PersonEmployerCompanyTrade      []string `json:"person_employer_company_trade"`
	PersonEmployerClientCompanyName string   `json:"person_employer_client_company_name"`
	PersonEmployerClientCompanyUEN  string   `json:"person_employer_client_company_unique_entity_number"`

	// Attendance
	PersonAttendanceDate    string             `json:"person_attendance_date"` // YYYY-MM-DD
	PersonAttendanceDetails []AttendanceDetail `json:"person_attendance_details"`
}

type AttendanceDetail struct {
	TimeIn  string `json:"time_in"`
	TimeOut string `json:"time_out"`
}

// ManpowerDistribution represents manpower distribution data for offsite fabrication
type ManpowerDistribution struct {
	SubmissionMonth                      string                       `json:"submission_month"` // format "YYYY-MM"
	OffsiteFabricatorCompanyName         string                       `json:"offsite_fabricator_company_name"`
	OffsiteFabricatorCompanyUEN          string                       `json:"offsite_fabricator_company_unique_entity_number"`
	OffsiteFabricatorLocationDescription string                       `json:"offsite_fabricator_location_description"`
	ManpowerDistributionStorageRatio     int                          `json:"manpower_distribution_storage_ratio"`
	ManpowerDistributionClientDetails    []ManpowerDistributionClient `json:"manpower_distribution_client_details"`
}

// ManpowerDistributionClient represents manpower allocation for a client project
type ManpowerDistributionClient struct {
	ProjectReferenceNumber     string `json:"project_reference_number"`
	ProjectTitle               string `json:"project_title"`
	ProjectLocationDescription string `json:"project_location_description"`
	FabricationStartMonth      string `json:"fabrication_start_month"`    // format "YYYY-MM"
	FabricationCompleteMonth   string `json:"fabrication_complete_month"` // format "YYYY-MM"
	ManpowerRatio              int    `json:"manpower_ratio"`             // e.g., 100 for 100%
}
