package payloads

// ManpowerUtilization represents the manpower utilization record for a project
type ManpowerUtilization struct {
	// Internal fields (not exported to JSON)
	InternalAttendanceID  string `json:"-"`
	InternalWorkerID      string `json:"-"`
	InternalSiteID        string `json:"-"`
	InternalRegulatorID   string `json:"-"`
	InternalRegulatorName string `json:"-"`
	InternalOnBehalfOfID  string `json:"-"`

	SubmissionEntity int    `json:"submission_entity"`
	SubmissionMonth  string `json:"submission_month"` // YYYY-MM

	// Onsite Builder (submission_entity = 1)
	ProjectReferenceNumber     *string `json:"project_reference_number,omitempty"`
	ProjectTitle               *string `json:"project_title,omitempty"`
	ProjectLocationDescription *string `json:"project_location_description,omitempty"`
	ProjectContractNumber      *string `json:"project_contract_number,omitempty"`
	ProjectContractName        *string `json:"project_contract_name,omitempty"`
	HdbPrecinctName            *string `json:"hdb_precinct_name,omitempty"`
	MainContractorCompanyName  *string `json:"main_contractor_company_name,omitempty"`
	MainContractorCompanyUEN   *string `json:"main_contractor_company_unique_entity_number,omitempty"`

	// Person
	PersonIDNo                      string   `json:"person_id_no"`
	PersonIDAndWorkPassType         string   `json:"person_id_and_work_pass_type"`
	PersonNationality               *string  `json:"person_nationality,omitempty"`
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
