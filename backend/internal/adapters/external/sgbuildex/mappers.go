package sgbuildex

import (
	"log"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/pkg/validation"
	"strings"
	"time"

	"sgbuildex/internal/adapters/external/sgbuildex/payloads"
)

// validateMandatoryFields enforces all API-mandatory fields before a row is submitted.
// It enforces the universal mandatory fields AND regulator-specific fields (BCA/HDB/LTA).
// Returns the name of the first missing field, or "" if all mandatory fields are present.
func validateMandatoryFields(r domain.AttendanceRow) string {
	// ── Universal mandatory (all regulators) ──
	universal := []struct {
		name string
		val  string
	}{
		{"person_id_no", r.WorkerFIN},
		{"person_id_and_work_pass_type", r.WorkerWorkPassType},
		{"person_trade", r.WorkerTrade},
		{"person_employer_company_name", r.EmployerName},
		{"person_employer_company_unique_entity_number", r.EmployerUEN},
	}
	for _, f := range universal {
		if strings.TrimSpace(f.val) == "" {
			return f.name
		}
	}

	// ── Regulator-specific mandatory fields ──
	reg := strings.ToUpper(strings.TrimSpace(r.RegulatorName))

	// BCA: person_employer_client_company_name, _uen, and person_employer_company_trade
	if reg == "BCA" {
		if strings.TrimSpace(r.EmployerClientName) == "" {
			return "person_employer_client_company_name (BCA mandatory)"
		}
		if strings.TrimSpace(r.EmployerClientUEN) == "" {
			return "person_employer_client_company_unique_entity_number (BCA mandatory)"
		}
		if strings.TrimSpace(r.EmployerTrade) == "" {
			return "person_employer_company_trade (BCA mandatory)"
		}
	}

	// LTA: person_employer_company_trade
	if reg == "LTA" {
		if strings.TrimSpace(r.EmployerTrade) == "" {
			return "person_employer_company_trade (LTA mandatory)"
		}
	}

	// HDB: person_nationality
	if reg == "HDB" {
		if strings.TrimSpace(r.WorkerNationality) == "" {
			return "person_nationality (HDB mandatory)"
		}
	}

	return "" // all checks passed
}

// MapAttendanceToManpower converts DB rows to ManpowerUtilization payloads.
// Records that are missing API-mandatory fields are skipped and logged.
func MapAttendanceToManpower(rows []domain.AttendanceRow) []payloads.ManpowerUtilization {
	var results []payloads.ManpowerUtilization
	for _, r := range rows {
		// Guard: skip rows missing any mandatory fields (universal + regulator-specific)
		if missing := validateMandatoryFields(r); missing != "" {
			log.Printf("[SGBuildex] SKIP attendance %s (regulator=%s worker=%s): mandatory field '%s' is empty",
				r.AttendanceID, r.RegulatorName, r.WorkerFIN, missing)
			continue
		}

		payload := payloads.ManpowerUtilization{
			InternalAttendanceID:            r.AttendanceID,
			InternalWorkerID:                r.WorkerID,
			InternalSiteID:                  r.SiteID,
			InternalRegulatorID:             r.RegulatorID,
			InternalRegulatorName:           r.RegulatorName,
			InternalOnBehalfOfID:            r.OnBehalfOfID,
			SubmissionEntity:                ptrIntOrDefault(r.SubmissionEntity, 1),
			SubmissionMonth:                 r.SubmissionDate.Format("2006-01"),
			PersonIDNo:                      Ptr(strings.ToUpper(strings.TrimSpace(r.WorkerFIN))),
			PersonIDAndWorkPassType:         Ptr(strings.ToUpper(strings.TrimSpace(r.WorkerWorkPassType))),
			PersonNationality:               Ptr(strings.ToUpper(strings.TrimSpace(r.WorkerNationality))),
			PersonTrade:                     Ptr(r.WorkerTrade),
			PersonEmployerCompanyName:       Ptr(r.EmployerName),
			PersonEmployerCompanyUEN:        Ptr(validation.SanitizeUEN(r.EmployerUEN)),
			PersonEmployerCompanyTrade:      parseTrades(r.EmployerTrade),
			PersonEmployerClientCompanyName: Ptr(r.EmployerClientName),
			PersonEmployerClientCompanyUEN:  Ptr(validation.SanitizeUEN(r.EmployerClientUEN)),
			PersonAttendanceDate:            r.TimeIn.Format("2006-01-02"),
			PersonAttendanceDetails: []payloads.AttendanceDetail{
				{
					TimeIn:  r.TimeIn.Format(time.RFC3339),
					TimeOut: FormatOptionalTime(r.TimeOut),
				},
			},
		}

		// Conditional fields based on Submission Entity
		if r.SubmissionEntity == 2 {
			// Offsite Fabricator (SubmissionEntity = 2)
			payload.OffsiteFabricatorCompanyName = Ptr(r.OffsiteFabricatorName)
			payload.OffsiteFabricatorCompanyUEN = Ptr(validation.SanitizeUEN(r.OffsiteFabricatorUEN))
			payload.OffsiteFabricatorLocationDescription = Ptr(r.OffsiteFabricatorLocation)
		} else {
			// Onsite Builder (SubmissionEntity = 1 - Default)
			payload.ProjectReferenceNumber = Ptr(r.ProjectRef)
			payload.ProjectTitle = Ptr(r.ProjectTitle)
			payload.ProjectLocationDescription = Ptr(r.ProjectLocation)
			payload.ProjectContractNumber = Ptr(r.ProjectContractNo)
			payload.ProjectContractName = Ptr(r.ProjectContractName)
			payload.HdbPrecinctName = Ptr(r.HDBPrecinctName)
			payload.MainContractorCompanyName = Ptr(r.SiteOwnerName)
			payload.MainContractorCompanyUEN = Ptr(validation.SanitizeUEN(r.SiteOwnerUEN))
		}

		results = append(results, payload)
	}

	return results
}

func parseTrades(tradeStr string) []string {
	trimmed := strings.TrimSpace(tradeStr)
	if trimmed == "" || strings.ToLower(trimmed) == "null" {
		return nil
	}
	parts := strings.Split(tradeStr, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
