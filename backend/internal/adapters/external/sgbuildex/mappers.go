package sgbuildex

import (
	"database/sql"
	"sgbuildex/internal/core/domain"
	"strings"
	"time"

	"sgbuildex/internal/adapters/external/sgbuildex/payloads"
)

// MapAttendanceToManpower converts DB rows to ManpowerUtilization payloads
func MapAttendanceToManpower(rows []domain.AttendanceRow) []payloads.ManpowerUtilization {
	ptr := func(s string) *string { return &s }

	var results []payloads.ManpowerUtilization
	for _, r := range rows {
		payload := payloads.ManpowerUtilization{
			InternalAttendanceID:            r.AttendanceID,
			InternalWorkerID:                r.WorkerID,
			InternalSiteID:                  r.SiteID,
			InternalPICName:                 r.PICName,
			InternalPICFIN:                  r.PICFIN,
			SubmissionEntity:                1, // 1 for Onsite Builder
			SubmissionMonth:                 r.SubmissionDate.Format("2006-01"),
			ProjectReferenceNumber:          ptr(r.ProjectRef),
			ProjectTitle:                    ptr(r.SiteName),
			ProjectLocationDescription:      ptr(r.SiteLocation),
			MainContractorCompanyName:       ptr(r.SiteOwnerName),
			MainContractorCompanyUEN:        ptr(r.SiteOwnerUEN),
			PersonIDNo:                      r.WorkerFIN,
			PersonName:                      r.WorkerName,
			PersonIDAndWorkPassType:         "NRIC", // Defaulting to NRIC, can be refined
			PersonTrade:                     r.TradeCode,
			PersonEmployerCompanyName:       r.EmployerName,
			PersonEmployerCompanyUEN:        r.EmployerUEN,
			PersonEmployerCompanyTrade:      parseTrades(r.EmployerTrade),
			PersonEmployerClientCompanyName: r.SiteOwnerName,
			PersonEmployerClientCompanyUEN:  r.SiteOwnerUEN,
			PersonAttendanceDate:            r.TimeIn.Format("2006-01-02"),
			PersonAttendanceDetails: []payloads.AttendanceDetail{
				{
					TimeIn:  r.TimeIn.Format(time.RFC3339),
					TimeOut: formatNullTime(r.TimeOut),
				},
			},
		}

		results = append(results, payload)
	}

	return results
}

func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(time.RFC3339)
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
