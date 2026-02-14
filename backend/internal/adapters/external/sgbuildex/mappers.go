package sgbuildex

import (
	"database/sql"
	"time"

	"sgbuildex/internal/adapters/external/sgbuildex/payloads"
)

// MapAttendanceToManpower converts DB rows to ManpowerUtilization payloads
func MapAttendanceToManpower(rows []AttendanceRow) []payloads.ManpowerUtilization {
	ptr := func(s string) *string { return &s }

	var results []payloads.ManpowerUtilization
	for _, r := range rows {
		payload := payloads.ManpowerUtilization{
			InternalAttendanceID:            r.AttendanceID,
			InternalWorkerID:                r.WorkerID,
			InternalSiteID:                  r.SiteID,
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
			PersonEmployerCompanyTrade:      []string{r.TradeCode},
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

// MapAggregationToDistribution converts aggregated SQL results to ManpowerDistribution payloads
func MapAggregationToDistribution(rows []MonthlyDistributionRow) []payloads.ManpowerDistribution {
	type fabricatorKey struct {
		UEN   string
		Month string
	}

	type fabricatorData struct {
		Name     string
		Location string
		Projects []payloads.ManpowerDistributionClient
		Total    int
	}

	data := make(map[fabricatorKey]*fabricatorData)

	// Group results by fabricator to calculate ratios
	for _, r := range rows {
		key := fabricatorKey{UEN: r.FabricatorUEN, Month: r.SubmissionMonth}
		if _, ok := data[key]; !ok {
			data[key] = &fabricatorData{
				Name:     r.FabricatorName,
				Location: r.FabricatorLocation,
				Projects: []payloads.ManpowerDistributionClient{},
			}
		}
		fData := data[key]
		fData.Total += r.AttendanceCount
		fData.Projects = append(fData.Projects, payloads.ManpowerDistributionClient{
			ProjectReferenceNumber:     r.ProjectRef,
			ProjectTitle:               r.ProjectTitle,
			ProjectLocationDescription: r.ProjectLocation,
			FabricationStartMonth:      r.SubmissionMonth,
			FabricationCompleteMonth:   r.SubmissionMonth,
			ManpowerRatio:              r.AttendanceCount, // Store count temporarily, will convert to ratio next
		})
	}

	var results []payloads.ManpowerDistribution
	for key, fData := range data {
		payload := payloads.ManpowerDistribution{
			SubmissionMonth:                      key.Month,
			OffsiteFabricatorCompanyName:         fData.Name,
			OffsiteFabricatorCompanyUEN:          key.UEN,
			OffsiteFabricatorLocationDescription: fData.Location,
			ManpowerDistributionStorageRatio:     20, // Default value
			ManpowerDistributionClientDetails:    []payloads.ManpowerDistributionClient{},
		}

		for _, p := range fData.Projects {
			ratio := 0
			if fData.Total > 0 {
				ratio = (p.ManpowerRatio * 100) / fData.Total
			}
			p.ManpowerRatio = ratio
			payload.ManpowerDistributionClientDetails = append(payload.ManpowerDistributionClientDetails, p)
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
