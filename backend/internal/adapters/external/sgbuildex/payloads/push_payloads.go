package payloads

import "time"

// LiftDatahubPayload represents the Lift Datahub event submission
type LiftDatahubPayload struct {
	ReportIdentificationNumber  *string    `json:"report_identification_number"`   // required
	PtoLiftIdentificationNumber *string    `json:"pto_lift_identification_number"` // required
	EventStartDateTime          time.Time  `json:"event_start_date_time"`          // required
	EventType                   int        `json:"event_type"`                     // required, enum 1..6
	EventEndDateTime            *time.Time `json:"event_end_date_time,omitempty"`
	EventLiftSystem             *int       `json:"event_lift_system,omitempty"` // enum 1..18
	EventDescription            *string    `json:"event_description,omitempty"`
	ActionType                  *int       `json:"action_type,omitempty"` // enum 1..6
	ActionDescription           *string    `json:"action_description,omitempty"`
	LiftPartReplacedDescription *string    `json:"lift_part_replaced_description,omitempty"`
	PersonFullName              *string    `json:"person_full_name,omitempty"` // max 66 chars
	VisitType                   *int       `json:"visit_type,omitempty"`       // enum 1..3
	FaultType                   *int       `json:"fault_type,omitempty"`       // enum 1..3
	Remarks                     *string    `json:"remarks,omitempty"`
}

// UltimateLoadTestPayload represents the Ultimate Load Test (ULT) data element
type UltimateLoadTestPayload struct {
	ProjectReferenceNumber                                         string     `json:"project_reference_number"`                                             // required
	ProjectTitle                                                   string     `json:"project_title"`                                                        // required
	ProjectLocationDescription                                     string     `json:"project_location_description"`                                         // required
	ProjectMainContractorCompanyName                               string     `json:"project_main_contractor_company_name"`                                 // required
	ProjectMainContractorCompanyUniqueEntityNumber                 string     `json:"project_main_contractor_company_unique_entity_number"`                 // required
	TechnicalControllerPersonName                                  string     `json:"technical_controller_person_name"`                                     // required
	RegisteredEngineerRegisteredTechnicalOfficerPersonName         string     `json:"registered_engineer_registered_technical_officer_person_name"`         // required
	RegisteredEngineerRegisteredTechnicalOfficerRegistrationNumber string     `json:"registered_engineer_registered_technical_officer_registration_number"` // required
	QualifiedPersonSupervisionPersonName                           string     `json:"qualified_person_supervision_person_name"`                             // required
	QualifiedPersonSupervisionRegistrationNumber                   string     `json:"qualified_person_supervision_registration_number"`                     // required
	QualifiedPersonGeotechnicalPersonName                          *string    `json:"qualified_person_geotechnical_person_name,omitempty"`
	QualifiedPersonGeotechnicalRegistrationNumber                  *string    `json:"qualified_person_geotechnical_registration_number,omitempty"`
	ProjectTotalUltimateLoadTest                                   int        `json:"project_total_ultimate_load_test"` // required
	StructuralPlanNumber                                           string     `json:"structural_plan_number"`           // required
	PileReferenceNumber                                            string     `json:"pile_reference_number"`            // required
	UltimateLoadTestDate                                           time.Time  `json:"ultimate_load_test_date"`          // required
	UltimateLoadTestMethod                                         int        `json:"ultimate_load_test_method"`        // required, enum 1..4
	UltimateLoadTestMethodOther                                    *string    `json:"ultimate_load_test_method_other,omitempty"`
	PileDiameter                                                   int        `json:"pile_diameter"`                              // required
	PileAsBuiltLength                                              float64    `json:"pile_as_built_length"`                       // required
	PileWorkingLoad                                                int        `json:"pile_working_load"`                          // required
	PileHeadSettlement15TimeWorkingLoad                            float64    `json:"pile_head_settlement_1_5_time_working_load"` // required
	PileHeadSettlement20TimeWorkingLoad                            float64    `json:"pile_head_settlement_2_0_time_working_load"` // required
	MaximumXTimeWorkingLoadBeforeFailure                           float64    `json:"maximum_x_time_working_load_before_failure"` // required
	PileHeadSettlementXTimeWorkingLoad                             float64    `json:"pile_head_settlement_x_time_working_load"`   // required
	UltimateLoadTestResult                                         int        `json:"ultimate_load_test_result"`                  // required, enum 1..4
	RedoneUltimateLoadTest                                         bool       `json:"redone_ultimate_load_test"`                  // required
	UltimateLoadTestAmendmentRemarks                               *string    `json:"ultimate_load_test_amendment_remarks,omitempty"`
	UltimateLoadTestAmendmentDate                                  *time.Time `json:"ultimate_load_test_amendment_date,omitempty"`
}

// WorkingLoadTestPayload

type WorkingLoadTestPayload struct {
	ProjectReferenceNumber                                         string     `json:"project_reference_number"`
	ProjectTitle                                                   string     `json:"project_title"`
	ProjectLocationDescription                                     string     `json:"project_location_description"`
	ProjectMainContractorCompanyName                               string     `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUniqueEntityNumber                 string     `json:"project_main_contractor_company_unique_entity_number"`
	TechnicalControllerPersonName                                  string     `json:"technical_controller_person_name"`
	RegisteredEngineerRegisteredTechnicalOfficerPersonName         string     `json:"registered_engineer_registered_technical_officer_person_name"`
	RegisteredEngineerRegisteredTechnicalOfficerRegistrationNumber string     `json:"registered_engineer_registered_technical_officer_registration_number"`
	QualifiedPersonSupervisionPersonName                           string     `json:"qualified_person_supervision_person_name"`
	QualifiedPersonSupervisionRegistrationNumber                   string     `json:"qualified_person_supervision_registration_number"`
	QualifiedPersonGeotechnicalPersonName                          *string    `json:"qualified_person_geotechnical_person_name,omitempty"`
	QualifiedPersonGeotechnicalRegistrationNumber                  *string    `json:"qualified_person_geotechnical_registration_number,omitempty"`
	ProjectTotalWorkingLoadTest                                    int        `json:"project_total_working_load_test"`
	StructuralPlanNumber                                           string     `json:"structural_plan_number"`
	PileReferenceNumber                                            string     `json:"pile_reference_number"`
	WorkingLoadTestDate                                            time.Time  `json:"working_load_test_date"`
	WorkingLoadTestMethod                                          int        `json:"working_load_test_method"`
	WorkingLoadTestMethodOther                                     *string    `json:"working_load_test_method_other,omitempty"`
	PileDiameter                                                   int        `json:"pile_diameter"`
	PileAsBuiltLength                                              float64    `json:"pile_as_built_length"`
	PileWorkingLoad                                                int        `json:"pile_working_load"`
	PileHeadSettlement15TimeWorkingLoad                            float64    `json:"pile_head_settlement_1_5_time_working_load"`
	PileHeadSettlement20TimeWorkingLoad                            float64    `json:"pile_head_settlement_2_0_time_working_load"`
	WorkingLoadTestResult                                          int        `json:"working_load_test_result"`
	AdditionalWorkingLoadTest                                      bool       `json:"additional_working_load_test"`
	WorkingLoadTestAmendmentRemarks                                *string    `json:"working_load_test_amendment_remarks,omitempty"`
	WorkingLoadTestAmendmentDate                                   *time.Time `json:"working_load_test_amendment_date,omitempty"`
}

// ----------------------
// Piling Installation Record
// ----------------------
type PilingInstallationRecordPayload struct {
	ProjectReferenceNumber                                         string     `json:"project_reference_number"`
	ProjectTitle                                                   string     `json:"project_title"`
	ProjectLocationDescription                                     string     `json:"project_location_description"`
	ProjectMainContractorCompanyName                               string     `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUniqueEntityNumber                 string     `json:"project_main_contractor_company_unique_entity_number"`
	ProjectPilingContractorCompanyName                             string     `json:"project_piling_contractor_company_name"`
	ProjectPilingContractorCompanyUniqueEntityNumber               string     `json:"project_piling_contractor_company_unique_entity_number"`
	ProjectLandSurveyorCompanyName                                 *string    `json:"project_land_surveyor_company_name,omitempty"`
	ProjectLandSurveyorCompanyUniqueEntityNumber                   *string    `json:"project_land_surveyor_company_unique_entity_number,omitempty"`
	TechnicalControllerPersonName                                  string     `json:"technical_controller_person_name"`
	RegisteredEngineerRegisteredTechnicalOfficerPersonName         string     `json:"registered_engineer_registered_technical_officer_person_name"`
	RegisteredEngineerRegisteredTechnicalOfficerRegistrationNumber string     `json:"registered_engineer_registered_technical_officer_registration_number"`
	QualifiedPersonSupervisionPersonName                           string     `json:"qualified_person_supervision_person_name"`
	QualifiedPersonSupervisionRegistrationNumber                   string     `json:"qualified_person_supervision_registration_number"`
	QualifiedPersonGeotechnicalPersonName                          *string    `json:"qualified_person_geotechnical_person_name,omitempty"`
	QualifiedPersonGeotechnicalRegistrationNumber                  *string    `json:"qualified_person_geotechnical_registration_number,omitempty"`
	LandSurveyorPersonName                                         *string    `json:"land_surveyor_person_name,omitempty"`
	LandSurveyorRegistrationNumber                                 *int64     `json:"land_surveyor_registration_number,omitempty"`
	ProjectTotalPiles                                              int64      `json:"project_total_piles"`
	StructuralPlanNumber                                           string     `json:"structural_plan_number"`
	PileReferenceNumber                                            string     `json:"pile_reference_number"`
	PilingInstallationDate                                         time.Time  `json:"piling_installation_date"`
	ProjectPilingWorkType                                          int        `json:"project_piling_work_type"`
	ProjectPilingWorkTypeOther                                     *string    `json:"project_piling_work_type_other,omitempty"`
	ProjectPilingFoundationTypeOther                               *string    `json:"project_piling_foundation_type_other,omitempty"`
	PileXEasting                                                   float64    `json:"pile_x_easting"`
	PileYNorthing                                                  float64    `json:"pile_y_northing"`
	PileCutOffLevel                                                float64    `json:"pile_cut_off_level"`
	PileToeLevel                                                   float64    `json:"pile_toe_level"`
	PileDiameterLongestLength                                      int        `json:"pile_diameter_longest_length"`
	PileWidth                                                      *int64     `json:"pile_width,omitempty"`
	PileDesignPenetrationLength                                    float64    `json:"pile_design_penetration_length"`
	PileAsBuiltLength                                              float64    `json:"pile_as_built_length"`
	PileDesignSocketingLength                                      *float64   `json:"pile_design_socketing_length,omitempty"`
	PileActualSocketingLength                                      *float64   `json:"pile_actual_socketing_length,omitempty"`
	PileDesignEmbedmentLength                                      *float64   `json:"pile_design_embedment_length,omitempty"`
	PileActualEmbedmentLength                                      *float64   `json:"pile_actual_embedment_length,omitempty"`
	PileLocalXEccentricity                                         *int64     `json:"pile_local_x_eccentricity,omitempty"`
	PileLocalYEccentricity                                         *int64     `json:"pile_local_y_eccentricity,omitempty"`
	PileZoneBoreHoleNumber                                         *string    `json:"pile_zone_bore_hole_number,omitempty"`
	PileBoringStartDate                                            *time.Time `json:"pile_boring_start_date,omitempty"`
	PileBoringCompleteDate                                         *time.Time `json:"pile_boring_complete_date,omitempty"`
	PileVerticality                                                *int64     `json:"pile_verticality,omitempty"`
	PileReinforcementBarNumberSize                                 *string    `json:"pile_reinforcement_bar_number_size,omitempty"`
	PileReinforcementBarLength                                     *float64   `json:"pile_reinforcement_bar_length,omitempty"`
	PileReinforcementLinkSizeSpacing                               *string    `json:"pile_reinforcement_link_size_spacing,omitempty"`
	PileSpacerSize                                                 *int64     `json:"pile_spacer_size,omitempty"`
	PileSpacerSpacing                                              *int64     `json:"pile_spacer_spacing,omitempty"`
	PileConcretingMethod                                           *int64     `json:"pile_concreting_method,omitempty"`
	PileConcretingMethodOther                                      *string    `json:"pile_concreting_method_other,omitempty"`
	PileConcreteGrade                                              string     `json:"pile_concrete_grade"`
	PileToeCleaned                                                 *bool      `json:"pile_toe_cleaned,omitempty"`
	PileConcretingStartDateTime                                    *time.Time `json:"pile_concreting_start_date_time,omitempty"`
	PileConcretingCompleteDateTime                                 *time.Time `json:"pile_concreting_complete_date_time,omitempty"`
	PileCalculatedConcreteVolume                                   *float64   `json:"pile_calculated_concrete_volume,omitempty"`
	PileActualConcreteVolume                                       *float64   `json:"pile_actual_concrete_volume,omitempty"`
	ProjectPileConcreteSupplierCompanyName                         *string    `json:"project_pile_concrete_supplier_company_name,omitempty"`
	CompetentSoilStandardPenetrationTestRequirement                *int64     `json:"competent_soil_standard_penetration_test_requirement,omitempty"`
	CompetentSoilDepth                                             *float64   `json:"competent_soil_depth,omitempty"`
}

// ConcreteCubeTestContractor represents a contractor's concrete cube test report
type ConcreteCubeTestContractor struct {
	ProjectReferenceNumber                       string     `json:"project_reference_number"`
	ProjectTitle                                 string     `json:"project_title"`
	ProjectLocationDescription                   string     `json:"project_location_description"`
	ProjectMainContractorCompanyName             string     `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUEN              string     `json:"project_main_contractor_company_unique_entity_number"`
	ProjectTestLaboratoryCompanyName             string     `json:"project_test_laboratory_company_name"`
	ProjectTestLaboratoryCompanyUEN              string     `json:"project_test_laboratory_company_unique_entity_number"`
	ProjectTestLaboratoryContractNumber          *string    `json:"project_test_laboratory_contract_number,omitempty"`
	ProjectConcreteSupplierCompanyName           *string    `json:"project_concrete_supplier_company_name,omitempty"`
	ProjectConcreteSupplierCompanyRegistrationNo *string    `json:"project_concrete_supplier_company_registration_number,omitempty"`
	RegisteredEngineerRTOPersonName              string     `json:"registered_engineer_registered_technical_officer_person_name"`
	RegisteredEngineerRTORegistrationNumber      string     `json:"registered_engineer_registered_technical_officer_registration_number"`
	CastingLocation                              string     `json:"casting_location"`
	CastingDate                                  time.Time  `json:"casting_date"`
	ConcreteSlump                                *int       `json:"concrete_slump,omitempty"`
	ConcreteGrade                                string     `json:"concrete_grade"`
	ConcreteType                                 int        `json:"concrete_type"`
	ConcreteTypeAdditionalInformation            *string    `json:"concrete_type_additional_information,omitempty"`
	ConcreteWorkType                             int        `json:"concrete_work_type"`
	ConcreteMixType                              *int       `json:"concrete_mix_type,omitempty"`
	ConcreteCubeTestJobReferenceNumber           string     `json:"concrete_cube_test_job_reference_number"`
	ConcreteCubeTestTypeRequired                 int        `json:"concrete_cube_test_type_required"`
	ConcreteCubeTestTypeOther                    *int       `json:"concrete_cube_test_type_other,omitempty"`
	ConcreteCubeTestDate                         time.Time  `json:"concrete_cube_test_date"`
	ConcreteCubeTestSampleAverageCubeStrength    float64    `json:"concrete_cube_test_sample_average_cube_strength"`
	ConcreteCubeTestRollingAverageCubeStrength   *float64   `json:"concrete_cube_test_rolling_average_cube_strength,omitempty"`
	ConcreteCubeTestResult                       bool       `json:"concrete_cube_test_result"`
	ConcreteCubeSizeAndTestStandard              *int       `json:"concrete_cube_size_and_test_standard,omitempty"`
	ConcreteCubeTestRectificationRemarks         *string    `json:"concrete_cube_test_rectification_remarks,omitempty"`
	ConcreteCubeTestRectificationDate            *time.Time `json:"concrete_cube_test_rectification_date,omitempty"`
	ConcreteCubeTestSampleNumber                 *string    `json:"concrete_cube_test_sample_number,omitempty"`

	ConcreteCubeTestAttachments *struct {
		Attachments []struct {
			Filename    string `json:"filename"`
			FileContent string `json:"file_content"`
		} `json:"attachments"`
	} `json:"concrete_cube_test_attachments,omitempty"`

	ConcreteCubeTestDetails []struct {
		CubeReferenceNumber string  `json:"cube_reference_number"`
		CubeMass            float64 `json:"cube_mass"`
		CubeDensity         float64 `json:"cube_density"`
		CubeFractureType    bool    `json:"cube_fracture_type"`
		CubeStrength        float64 `json:"cube_strength"`
	} `json:"concrete_cube_test_details"`
}

// ConcreteCubeTestContractor represents a contractor's concrete cube test report
type ConcreteCubeTestLaboratory struct {
	ProjectReferenceNumber                       string     `json:"project_reference_number"`
	ProjectTitle                                 string     `json:"project_title"`
	ProjectLocationDescription                   string     `json:"project_location_description"`
	ProjectMainContractorCompanyName             string     `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUEN              string     `json:"project_main_contractor_company_unique_entity_number"`
	ProjectTestLaboratoryCompanyName             string     `json:"project_test_laboratory_company_name"`
	ProjectTestLaboratoryCompanyUEN              string     `json:"project_test_laboratory_company_unique_entity_number"`
	ProjectTestLaboratoryContractNumber          *string    `json:"project_test_laboratory_contract_number,omitempty"`
	ProjectConcreteSupplierCompanyName           *string    `json:"project_concrete_supplier_company_name,omitempty"`
	ProjectConcreteSupplierCompanyRegistrationNo *string    `json:"project_concrete_supplier_company_registration_number,omitempty"`
	RegisteredEngineerRTOPersonName              string     `json:"registered_engineer_registered_technical_officer_person_name"`
	RegisteredEngineerRTORegistrationNumber      string     `json:"registered_engineer_registered_technical_officer_registration_number"`
	CastingLocation                              string     `json:"casting_location"`
	CastingDate                                  time.Time  `json:"casting_date"`
	ConcreteSlump                                *int       `json:"concrete_slump,omitempty"`
	ConcreteGrade                                string     `json:"concrete_grade"`
	ConcreteType                                 int        `json:"concrete_type"`
	ConcreteTypeAdditionalInformation            *string    `json:"concrete_type_additional_information,omitempty"`
	ConcreteWorkType                             int        `json:"concrete_work_type"`
	ConcreteMixType                              *int       `json:"concrete_mix_type,omitempty"`
	ConcreteCubeTestJobReferenceNumber           string     `json:"concrete_cube_test_job_reference_number"`
	ConcreteCubeTestTypeRequired                 int        `json:"concrete_cube_test_type_required"`
	ConcreteCubeTestTypeOther                    *int       `json:"concrete_cube_test_type_other,omitempty"`
	ConcreteCubeTestDate                         time.Time  `json:"concrete_cube_test_date"`
	ConcreteCubeTestSampleAverageCubeStrength    float64    `json:"concrete_cube_test_sample_average_cube_strength"`
	ConcreteCubeTestRollingAverageCubeStrength   *float64   `json:"concrete_cube_test_rolling_average_cube_strength,omitempty"`
	ConcreteCubeTestResult                       bool       `json:"concrete_cube_test_result"`
	ConcreteCubeSizeAndTestStandard              *int       `json:"concrete_cube_size_and_test_standard,omitempty"`
	ConcreteCubeTestRectificationRemarks         *string    `json:"concrete_cube_test_rectification_remarks,omitempty"`
	ConcreteCubeTestRectificationDate            *time.Time `json:"concrete_cube_test_rectification_date,omitempty"`
	ConcreteCubeTestSampleNumber                 *string    `json:"concrete_cube_test_sample_number,omitempty"`

	ConcreteCubeTestAttachments *struct {
		Attachments []struct {
			Filename    string `json:"filename"`
			FileContent string `json:"file_content"`
		} `json:"attachments"`
	} `json:"concrete_cube_test_attachments,omitempty"`

	ConcreteCubeTestDetails []struct {
		CubeReferenceNumber string  `json:"cube_reference_number"`
		CubeMass            float64 `json:"cube_mass"`
		CubeDensity         float64 `json:"cube_density"`
		CubeFractureType    bool    `json:"cube_fracture_type"`
		CubeStrength        float64 `json:"cube_strength"`
	} `json:"concrete_cube_test_details"`
}

// SteelElementTest represents a steel element test report
type SteelElementTest struct {
	ProjectReferenceNumber                          string `json:"project_reference_number"`
	ProjectTitle                                    string `json:"project_title"`
	ProjectLocationDescription                      string `json:"project_location_description"`
	ProjectMainContractorCompanyName                string `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUEN                 string `json:"project_main_contractor_company_unique_entity_number"`
	ProjectTestLaboratoryCompanyName                string `json:"project_test_laboratory_company_name"`
	ProjectTestLaboratoryCompanyUEN                 string `json:"project_test_laboratory_company_unique_entity_number"`
	ProjectSteelMillCompanyName                     string `json:"project_steel_mill_company_name"`
	ProjectSteelMillCompanyCountryManufacture       string `json:"project_steel_mill_company_country_manufacture"`
	ProjectSteelFabricatorCompanyName               string `json:"project_steel_fabricator_company_name"`
	ProjectSteelFabricatorCompanyCountryFabrication string `json:"project_steel_fabricator_company_country_fabrication"`
	ProjectInspectionTestingAgencyCompanyName       string `json:"project_inspection_testing_agency_company_name"`
	ProjectInspectionTestingAgencyCompanyUEN        string `json:"project_inspection_testing_agency_company_unique_entity_number"`
	SteelFabricatorAccreditationBLS                 bool   `json:"steel_fabricator_accreditation_builders_licensing_scheme"`
	TestLaboratoryAccreditationSAC                  bool   `json:"test_laboratory_accreditation_singapore_accreditation_council"`
	RegisteredEngineerRTOPersonName                 string `json:"registered_engineer_registered_technical_officer_person_name"`
	RegisteredEngineerRTORegistrationNumber         string `json:"registered_engineer_registered_technical_officer_registration_number"`

	SteelGrade          string  `json:"steel_grade"`
	SteelDesignStrength float64 `json:"steel_design_strength"`
	BoltSpecification   string  `json:"bolt_specification"`

	SteelElementTestReportNumber string    `json:"steel_element_test_report_number"`
	SteelElementTestDate         time.Time `json:"steel_element_test_date"`
	SteelElementTestResult       bool      `json:"steel_element_test_result"`
	SteelElementTestStandard     string    `json:"steel_element_test_standard"`
	SteelElementTestRemarks      string    `json:"steel_element_test_remarks"`

	SteelElementTestAttachments *struct {
		Attachments []struct {
			Filename    string `json:"filename"`
			FileContent string `json:"file_content"`
		} `json:"attachments"`
	} `json:"steel_element_test_attachments,omitempty"`

	SteelElementTestDetails []struct {
		SampleNumber    string  `json:"sample_number"`
		YieldStrength   float64 `json:"yield_strength"`
		TensileStrength float64 `json:"tensile_strength"`
		Elongation      float64 `json:"elongation"`
	} `json:"steel_element_test_details"`
}

// SteelRebarTest represents a steel rebar test report
type SteelRebarTest struct {
	ProjectReferenceNumber                          string `json:"project_reference_number"`
	ProjectTitle                                    string `json:"project_title"`
	ProjectLocationDescription                      string `json:"project_location_description"`
	ProjectMainContractorCompanyName                string `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUEN                 string `json:"project_main_contractor_company_unique_entity_number"`
	ProjectTestLaboratoryCompanyName                string `json:"project_test_laboratory_company_name"`
	ProjectTestLaboratoryCompanyUEN                 string `json:"project_test_laboratory_company_unique_entity_number"`
	ProjectSteelMillCompanyName                     string `json:"project_steel_mill_company_name"`
	ProjectSteelMillCompanyCountryManufacture       string `json:"project_steel_mill_company_country_manufacture"`
	ProjectSteelFabricatorCompanyName               string `json:"project_steel_fabricator_company_name"`
	ProjectSteelFabricatorCompanyCountryFabrication string `json:"project_steel_fabricator_company_country_fabrication"`
	TestLaboratoryAccreditationSAC                  bool   `json:"test_laboratory_accreditation_singapore_accreditation_council"`
	RegisteredEngineerRTOPersonName                 string `json:"registered_engineer_registered_technical_officer_person_name"`
	RegisteredEngineerRTORegistrationNumber         string `json:"registered_engineer_registered_technical_officer_registration_number"`

	SteelGrade                 string  `json:"steel_grade"`
	SteelDesignStrengthMinimum float64 `json:"steel_design_strength_minimum"`
	SteelDesignStrengthMaximum float64 `json:"steel_design_strength_maximum"`

	SteelRebarTestReportNumber string    `json:"steel_rebar_test_report_number"`
	SteelRebarTestDate         time.Time `json:"steel_rebar_test_date"`
	SteelRebarTestResult       bool      `json:"steel_rebar_test_result"`
	SteelRebarTestStandard     string    `json:"steel_rebar_test_standard"`
	SteelRebarTestRemarks      string    `json:"steel_rebar_test_remarks"`

	SteelRebarTestAttachments *struct {
		Attachments []struct {
			Filename    string `json:"filename"`
			FileContent string `json:"file_content"`
		} `json:"attachments"`
	} `json:"steel_rebar_test_attachments,omitempty"`

	SteelRebarTestDetails []struct {
		SampleNumber              string  `json:"sample_number"`
		NominalSize               int     `json:"nominal_size"`
		Mass                      float64 `json:"mass"`
		MeasuredLength            float64 `json:"measured_length"`
		NominalCrossSectionalArea float64 `json:"nominal_cross_sectional_area"`
		YieldPointLoad            float64 `json:"yield_point_load"`
		TensileTestYieldStrength  float64 `json:"tensile_test_yield_strength"`
		BendTestResult            int     `json:"bend_test_result"`
		RebendTestResult          int     `json:"rebend_test_result"`
	} `json:"steel_rebar_test_details"`
}

// SiteProgress represents the site progress update for a project
type SiteProgress struct {
	ProjectReferenceNumber           string    `json:"project_reference_number"`
	ProjectTitle                     string    `json:"project_title"`
	ProjectLocationDescription       string    `json:"project_location_description"`
	ProjectMainContractorCompanyName string    `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUEN  string    `json:"project_main_contractor_company_unique_entity_number"`
	DateOfUpdate                     time.Time `json:"date_of_update"`
	TotalNumberOfBlockZone           int       `json:"total_number_of_block_zone"`
	BlockZoneName                    string    `json:"block_zone_name"`
	ProjectStatus                    int       `json:"project_status"`
	DemolitionProgress               int       `json:"demolition_progress"`
	ERSSProgress                     int       `json:"erss_progress"`
	PilingProgress                   int       `json:"piling_progress"`
	SubstructureProgress             int       `json:"substructure_progress"`
	SuperstructureProgress           int       `json:"superstructure_progress"`

	SiteProgressAttachments *struct {
		Attachments []struct {
			Filename    string `json:"filename"`
			FileContent string `json:"file_content"`
		} `json:"attachments"`
	} `json:"site_progress_attachments,omitempty"`
}

// QualifiedPersonAttendance represents a qualified person or site supervisor attendance record
type QualifiedPersonAttendance struct {
	ProjectReferenceNumber                                 string    `json:"project_reference_number"`
	ProjectTitle                                           string    `json:"project_title"`
	ProjectLocationDescription                             string    `json:"project_location_description"`
	ProjectMainContractorCompanyName                       string    `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUEN                        string    `json:"project_main_contractor_company_unique_entity_number"`
	QualifiedPersonSiteSupervisorPersonName                string    `json:"qualified_person_site_supervisor_person_name"`
	WorkCategory                                           int       `json:"work_category"`
	ProfessionalEngineerRegistrationNumber                 string    `json:"professional_engineer_registration_number"`
	RegisteredEngineerRegisteredTechnicalOfficerRegNo      string    `json:"registered_engineer_registered_technical_officer_registration_number"`
	RegisteredEngineerRegisteredTechnicalOfficerTypeOfWork int       `json:"registered_engineer_registered_technical_officer_type_of_work"`
	TimeIn                                                 time.Time `json:"time_in"`
	TimeOut                                                time.Time `json:"time_out"`
	PurposeToEnteringSite                                  int       `json:"purpose_to_entering_site"`
}

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

// ERSSApproval represents the ERSS approval and inspection data for a project
type ERSSApproval struct {
	ProjectReferenceNumber           string      `json:"project_reference_number"`
	ProjectTitle                     string      `json:"project_title"`
	ProjectLocationDescription       string      `json:"project_location_description"`
	ProjectMainContractorCompanyName string      `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUEN  string      `json:"project_main_contractor_company_unique_entity_number"`
	TotalNumberOfBlockZone           int         `json:"total_number_of_block_zone"`
	BlockZoneName                    string      `json:"block_zone_name"`
	GridlineNumber                   string      `json:"gridline_number"`
	PartialLocationAttachments       Attachments `json:"partial_location_attachments"`
	GeotechnicalBuildingWorks        bool        `json:"geotechnical_building_works"`
	DeclarationForAnnexC1Submission  bool        `json:"declaration_for_annex_c1_submission"`
	CriticalConstructionStage        int         `json:"critical_construction_stage"`
	CriticalConstructionStageOthers  string      `json:"critical_construction_stage_others"`
	InstalledOrRemovedStrutSupport   int         `json:"installed_or_removed_strut_support"`
	TotalStrutSupport                int         `json:"total_strut_support"`

	SectionATechnicalControllerDeclarationDate string `json:"section_a_technical_controller_declaration_date"` // format "YYYY-MM-DDTHH:mm:ssZ"
	TechnicalControllerPersonName              string `json:"technical_controller_person_name"`

	SectionB1QualifiedPersonSupervisionInspectionDate       string `json:"section_b1_qualified_person_supervision_inspection_date"` // format "YYYY-MM-DDTHH:mm:ssZ"
	QualifiedPersonSupervisionDeviationStatus               bool   `json:"qualified_person_supervision_deviation_status"`
	QualifiedPersonSupervisionDeviationComments             string `json:"qualified_person_supervision_deviation_comments"`
	QualifiedPersonSupervisionGeotechnicalInspectionDate    string `json:"qualified_person_supervision_geotechnical_inspection_date_for_geotechnical_building_works"` // format "YYYY-MM-DDTHH:mm:ssZ"
	QualifiedPersonSupervisionGeotechnicalDeviationStatus   bool   `json:"qualified_person_supervision_geotechnical_deviation_status"`
	QualifiedPersonSupervisionGeotechnicalDeviationComments string `json:"qualified_person_supervision_geotechnical_deviation_comments"`

	SectionB2QualifiedPersonSupervisionApprovalDate string `json:"section_b2_qualified_person_supervision_approval_date"` // format "YYYY-MM-DDTHH:mm:ssZ"
	QualifiedPersonSupervisionAssessment            bool   `json:"qualified_person_supervision_assessment"`
	QualifiedPersonSupervisionName                  string `json:"qualified_person_supervision_name"`
	QualifiedPersonSupervisionRegistrationNumber    string `json:"qualified_person_supervision_registration_number"`

	QualifiedPersonSupervisionGeotechnicalApprovalDate              string `json:"qualified_person_supervision_geotechnical_approval_date_for_geotechnical_building_works"` // format "YYYY-MM-DDTHH:mm:ssZ"
	QualifiedPersonSupervisionGeotechnicalName                      string `json:"qualified_person_supervision_geotechnical_name"`
	QualifiedPersonSupervisionGeotechnicalProfessionalEngineerRegNo string `json:"qualified_person_supervision_geotechnical_professional_engineer_registration_number"`
}

// Attachments represents a generic attachment object
type Attachments struct {
	Attachments []Attachment `json:"attachments"`
}

// BuildingSettlementMonitoring represents monitoring of building settlement data
type BuildingSettlementMonitoring struct {
	ProjectReferenceNumber           string `json:"project_reference_number"`
	ProjectTitle                     string `json:"project_title"`
	ProjectLocationDescription       string `json:"project_location_description"`
	ProjectMainContractorCompanyName string `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUEN  string `json:"project_main_contractor_company_unique_entity_number"`

	TotalNumberOfBlockZone                  int    `json:"total_number_of_block_zone"`
	BlockZoneName                           string `json:"block_zone_name"`
	MaximumNumberOfStoreysExcludingBasement int    `json:"maximum_number_of_storeys_excluding_basement"`
	NumberOfStoreysReachedExcludingBasement int    `json:"number_of_storeys_reached_excluding_basement"`

	GeotechnicalBuildingWorks bool `json:"geotechnical_building_works"`

	QualifiedPersonSupervisionStructuralName                        string `json:"qualified_person_supervision_structural_name"`
	QualifiedPersonSupervisionStructuralRegistrationNumber          string `json:"qualified_person_supervision_structural_registration_number"`
	QualifiedPersonSupervisionGeotechnicalName                      string `json:"qualified_person_supervision_geotechnical_name"`
	QualifiedPersonSupervisionGeotechnicalProfessionalEngineerRegNo string `json:"qualified_person_supervision_geotechnical_professional_engineer_registration_number"`

	RecordingDateTime string `json:"recording_date_time"` // format "YYYY-MM-DDTHH:mm:ssZ"

	AllowableBuildingSettlement                   float64 `json:"allowable_building_settlement"`                     // in mm
	MaximumMeasuredBuildingSettlement             float64 `json:"maximum_measured_building_settlement"`              // in mm
	AllowableDifferentialBuildingSettlement       float64 `json:"allowable_differential_building_settlement"`        // in mm
	MaximumMeasuredDifferentialBuildingSettlement float64 `json:"maximum_measured_differential_building_settlement"` // in mm

	QualifiedPersonSupervisionObservation   bool   `json:"qualified_person_supervision_observation"`
	QualifiedPersonSupervisionOtherComments string `json:"qualified_person_supervision_other_comments"`
}

// NotificationToCBC represents a notification to the Commissioner of Building Control
type NotificationToCBC struct {
	ProjectReferenceNumber           string `json:"project_reference_number"`
	ProjectTitle                     string `json:"project_title"`
	ProjectLocationDescription       string `json:"project_location_description"`
	ProjectProcessingOfficerEmail    string `json:"project_processing_officer_email"`
	ProjectMainContractorCompanyName string `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUEN  string `json:"project_main_contractor_company_unique_entity_number"`

	RegisteredEngineerRegisteredTechnicalOfficerPersonName  string `json:"registered_engineer_registered_technical_officer_person_name"`
	RegisteredEngineerRegisteredTechnicalOfficerReRegNumber string `json:"registered_engineer_registered_technical_officer_re_registration_number"`

	QualifiedPersonSupervisionStructuralPersonName         string `json:"qualified_person_supervision_structural_person_name"`
	QualifiedPersonSupervisionStructuralRegistrationNumber string `json:"qualified_person_supervision_structural_registration_number"`

	QualifiedPersonDesignStructuralPersonName         string `json:"qualified_person_design_structural_person_name"`
	QualifiedPersonDesignStructuralRegistrationNumber string `json:"qualified_person_design_structural_registration_number"`

	NotificationCommissionerBuildingControlAttachments FileAttachments `json:"notification_commissioner_building_control_attachments"`

	ProjectNotificationDate        string `json:"project_notification_date"` // ISO 8601 format
	ProjectNotificationType        int    `json:"project_notification_type"`
	ProjectDocumentReferenceNumber string `json:"project_document_reference_number"`
	ProjectNotificationDescription string `json:"project_notification_description"`
	ProjectNotificationStatus      int    `json:"project_notification_status"`

	ProjectRectificationRemarks string `json:"project_rectification_remarks"`
	ProjectRectificationDate    string `json:"project_rectification_date"` // ISO 8601 format
}

// FileAttachments represents generic attachments for a project
type FileAttachments struct {
	Attachments []Attachment `json:"attachments"`
}

// Attachment represents a single attachment file
type Attachment struct {
	Filename    string `json:"filename"`
	FileContent string `json:"file_content"` // base64-encoded content
}

// ProjectDocument represents a project document submission
type ProjectDocument struct {
	ProjectReferenceNumber           string `json:"project_reference_number"`
	ProjectTitle                     string `json:"project_title"`
	ProjectLocationDescription       string `json:"project_location_description"`
	ProjectMainContractorCompanyName string `json:"project_main_contractor_company_name"`
	ProjectMainContractorCompanyUEN  string `json:"project_main_contractor_company_unique_entity_number"`

	RegisteredEngineerRegisteredTechnicalOfficerPersonName  string `json:"registered_engineer_registered_technical_officer_person_name"`
	RegisteredEngineerRegisteredTechnicalOfficerReRegNumber string `json:"registered_engineer_registered_technical_officer_re_registration_number"`

	QualifiedPersonSupervisionStructuralPersonName         string `json:"qualified_person_supervision_structural_person_name"`
	QualifiedPersonSupervisionStructuralRegistrationNumber string `json:"qualified_person_supervision_structural_registration_number"`

	QualifiedPersonSupervisionGeotechnicalPersonName         string `json:"qualified_person_supervision_geotechnical_person_name"`
	QualifiedPersonSupervisionGeotechnicalRegistrationNumber string `json:"qualified_person_supervision_geotechnical_registration_number"`

	QualifiedPersonDesignStructuralPersonName         string `json:"qualified_person_design_structural_person_name"`
	QualifiedPersonDesignStructuralRegistrationNumber string `json:"qualified_person_design_structural_registration_number"`

	ProjectDocumentType            int    `json:"project_document_type"` // e.g. 1-10
	ProjectDocumentTypeOther       string `json:"project_document_type_other"`
	ProjectDocumentReferenceNumber string `json:"project_document_reference_number"`
	ProjectDocumentRemarks         string `json:"project_document_remarks"`
	Capture360WebLink              string `json:"360_capture_web_link"`

	ProjectDocumentAttachments FileAttachments `json:"project_document_attachments"`
}
