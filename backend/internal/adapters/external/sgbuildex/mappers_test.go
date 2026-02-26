package sgbuildex

import (
	"sgbuildex/internal/core/domain"
	"testing"
)

func TestMapAggregationToDistribution(t *testing.T) {
	rows := []domain.MonthlyDistributionRow{
		{
			FabricatorName:     "Fab 1",
			FabricatorUEN:      "UEN1",
			FabricatorLocation: "Loc 1",
			ProjectRef:         "PRJ-A",
			ProjectTitle:       "Site A",
			ProjectLocation:    "Loc A",
			SubmissionMonth:    "2026-02",
			AttendanceCount:    10,
		},
		{
			FabricatorName:     "Fab 1",
			FabricatorUEN:      "UEN1",
			FabricatorLocation: "Loc 1",
			ProjectRef:         "PRJ-B",
			ProjectTitle:       "Site B",
			ProjectLocation:    "Loc B",
			SubmissionMonth:    "2026-02",
			AttendanceCount:    5,
		},
	}

	results := MapAggregationToDistribution(rows)

	if len(results) != 1 {
		t.Fatalf("Expected 1 fabricator result, got %d", len(results))
	}

	res := results[0]
	if res.OffsiteFabricatorCompanyUEN != "UEN1" {
		t.Errorf("Expected UEN1, got %s", res.OffsiteFabricatorCompanyUEN)
	}

	if len(res.ManpowerDistributionClientDetails) != 2 {
		t.Fatalf("Expected 2 project details, got %d", len(res.ManpowerDistributionClientDetails))
	}

	for _, detail := range res.ManpowerDistributionClientDetails {
		if detail.ProjectReferenceNumber == "PRJ-A" {
			if detail.ManpowerRatio != 66 { // (10/15)*100
				t.Errorf("Expected ratio 66 for PRJ-A, got %d", detail.ManpowerRatio)
			}
		} else if detail.ProjectReferenceNumber == "PRJ-B" {
			if detail.ManpowerRatio != 33 { // (5/15)*100
				t.Errorf("Expected ratio 33 for PRJ-B, got %d", detail.ManpowerRatio)
			}
		}
	}
}
