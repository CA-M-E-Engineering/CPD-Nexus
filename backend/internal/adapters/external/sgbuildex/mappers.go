package sgbuildex

import (
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/pkg/validation"
	"strings"
	"time"

	"sgbuildex/internal/adapters/external/sgbuildex/payloads"
)

// MapAttendanceToManpower converts DB rows to ManpowerUtilization payloads
func MapAttendanceToManpower(rows []domain.AttendanceRow) []payloads.ManpowerUtilization {
	var results []payloads.ManpowerUtilization
	for _, r := range rows {
		payload := payloads.ManpowerUtilization{
			InternalAttendanceID:            r.AttendanceID,
			InternalWorkerID:                r.WorkerID,
			InternalSiteID:                  r.SiteID,
			InternalRegulatorID:             r.RegulatorID,
			InternalRegulatorName:           r.RegulatorName,
			InternalOnBehalfOfID:            r.OnBehalfOfID,
			SubmissionEntity:                1, // 1 for Onsite Builder
			SubmissionMonth:                 r.SubmissionDate.Format("2006-01"),
			ProjectReferenceNumber:          Ptr(r.ProjectRef),
			ProjectTitle:                    Ptr(r.ProjectTitle),
			ProjectLocationDescription:      Ptr(r.ProjectLocation),
			ProjectContractNumber:           Ptr(r.ProjectContractNo),
			ProjectContractName:             Ptr(r.ProjectContractName),
			HdbPrecinctName:                 Ptr(r.HDBPrecinctName),
			MainContractorCompanyName:       Ptr(r.SiteOwnerName),
			MainContractorCompanyUEN:        Ptr(validation.SanitizeUEN(r.SiteOwnerUEN)),
			PersonIDNo:                      strings.ToUpper(strings.TrimSpace(r.WorkerFIN)),
			PersonIDAndWorkPassType:         strings.ToUpper(strings.TrimSpace(r.WorkerWorkPassType)),
			PersonNationality:               Ptr(strings.ToUpper(strings.TrimSpace(r.WorkerNationality))),
			PersonTrade:                     r.TradeCode,
			PersonEmployerCompanyName:       r.EmployerName,
			PersonEmployerCompanyUEN:        validation.SanitizeUEN(r.EmployerUEN),
			PersonEmployerCompanyTrade:      parseTrades(r.EmployerTrade),
			PersonEmployerClientCompanyName: r.EmployerClientName,
			PersonEmployerClientCompanyUEN:  validation.SanitizeUEN(r.EmployerClientUEN),
			PersonAttendanceDate:            r.TimeIn.Format("2006-01-02"),
			PersonAttendanceDetails: []payloads.AttendanceDetail{
				{
					TimeIn:  r.TimeIn.Format(time.RFC3339),
					TimeOut: FormatOptionalTime(r.TimeOut),
				},
			},
		}

		results = append(results, payload)
	}

	return results
}

func parseTrades(tradeStr string) []string {
	if tradeStr == "" {
		return []string{}
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
