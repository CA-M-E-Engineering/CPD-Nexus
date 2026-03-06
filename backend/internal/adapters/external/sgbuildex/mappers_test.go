package sgbuildex

import (
	"sgbuildex/internal/core/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateMandatoryFields(t *testing.T) {
	tests := []struct {
		name     string
		row      domain.AttendanceRow
		expected string // Expected missing field error message, empty if no error
	}{
		{
			name: "Valid row passes",
			row: domain.AttendanceRow{
				WorkerFIN:          "G1234567P",
				WorkerWorkPassType: "WP",
				WorkerTrade:        "2.3",
				EmployerName:       "Test Employer",
				EmployerUEN:        "12345678A",
				RegulatorID:        "REG-123",
				OnBehalfOfID:       "OB-456",
				RegulatorName:      "Other", // No specific regulatory checks
			},
			expected: "",
		},
		{
			name: "Missing universal field (WorkerFIN)",
			row: domain.AttendanceRow{
				WorkerWorkPassType: "WP",
				WorkerTrade:        "2.3",
				EmployerName:       "Test Employer",
				EmployerUEN:        "12345678A",
				RegulatorID:        "REG-123",
				OnBehalfOfID:       "OB-456",
			},
			expected: "person_id_no",
		},
		{
			name: "Missing infrastructure field (RegulatorID)",
			row: domain.AttendanceRow{
				WorkerFIN:          "G1234567P",
				WorkerWorkPassType: "WP",
				WorkerTrade:        "2.3",
				EmployerName:       "Test Employer",
				EmployerUEN:        "12345678A",
				OnBehalfOfID:       "OB-456",
			},
			expected: "regulator_id (Pitstop Configuration sync missing valid ID)",
		},
		{
			name: "Missing infrastructure field (OnBehalfOfID)",
			row: domain.AttendanceRow{
				WorkerFIN:          "G1234567P",
				WorkerWorkPassType: "WP",
				WorkerTrade:        "2.3",
				EmployerName:       "Test Employer",
				EmployerUEN:        "12345678A",
				RegulatorID:        "REG-123",
			},
			expected: "on_behalf_of_id (Pitstop Configuration sync missing valid UEN)",
		},
		{
			name: "BCA regulator missing client name",
			row: domain.AttendanceRow{
				WorkerFIN:          "G1234567P",
				WorkerWorkPassType: "WP",
				WorkerTrade:        "2.3",
				EmployerName:       "Test Employer",
				EmployerUEN:        "12345678A",
				RegulatorID:        "REG-123",
				OnBehalfOfID:       "OB-456",
				RegulatorName:      "BCA",
				EmployerClientUEN:  "87654321B",
				EmployerTrade:      "1.4",
			},
			expected: "person_employer_client_company_name (BCA mandatory)",
		},
		{
			name: "LTA regulator missing employer trade",
			row: domain.AttendanceRow{
				WorkerFIN:          "G1234567P",
				WorkerWorkPassType: "WP",
				WorkerTrade:        "2.3",
				EmployerName:       "Test Employer",
				EmployerUEN:        "12345678A",
				RegulatorID:        "REG-123",
				OnBehalfOfID:       "OB-456",
				RegulatorName:      "LTA",
			},
			expected: "person_employer_company_trade (LTA mandatory)",
		},
		{
			name: "HDB regulator missing nationality",
			row: domain.AttendanceRow{
				WorkerFIN:          "G1234567P",
				WorkerWorkPassType: "WP",
				WorkerTrade:        "2.3",
				EmployerName:       "Test Employer",
				EmployerUEN:        "12345678A",
				RegulatorID:        "REG-123",
				OnBehalfOfID:       "OB-456",
				RegulatorName:      "HDB",
			},
			expected: "person_nationality (HDB mandatory)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateMandatoryFields(tt.row)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapAttendanceToManpower(t *testing.T) {
	now := time.Now()

	validRow1 := domain.AttendanceRow{
		AttendanceID:       "ATT-1",
		WorkerID:           "W-1",
		SiteID:             "S-1",
		RegulatorID:        "REG-1",
		RegulatorName:      "BCA",
		OnBehalfOfID:       "OB-1",
		SubmissionEntity:   1, // Onsite Builder
		SubmissionDate:     now,
		WorkerFIN:          "G1234567P",
		WorkerWorkPassType: "WP",
		WorkerNationality:  "CN",
		WorkerTrade:        "2.3",
		EmployerName:       "Valid Employer",
		EmployerUEN:        "11111111A",
		EmployerTrade:      "1.1,2.2",
		EmployerClientName: "Valid Client",
		EmployerClientUEN:  "22222222B",
		TimeIn:             now,
		ProjectRef:         "REF-1",
		ProjectTitle:       "Title 1",
		ProjectLocation:    "Loc 1",
		SiteOwnerName:      "Main Con",
		SiteOwnerUEN:       "33333333C",
	}

	invalidRow := domain.AttendanceRow{
		AttendanceID: "ATT-2",
		WorkerFIN:    "", // Missing mandatory field
	}

	validRow2 := domain.AttendanceRow{
		AttendanceID:              "ATT-3",
		RegulatorID:               "REG-2",
		OnBehalfOfID:              "OB-2",
		SubmissionEntity:          2, // Offsite Fabricator
		WorkerFIN:                 "G7654321P",
		WorkerWorkPassType:        "WP",
		WorkerTrade:               "2.4",
		EmployerName:              "Valid Employer 2",
		EmployerUEN:               "44444444D",
		TimeIn:                    now,
		OffsiteFabricatorName:     "Fabricator",
		OffsiteFabricatorUEN:      "55555555E",
		OffsiteFabricatorLocation: "Fab Loc",
	}

	rows := []domain.AttendanceRow{validRow1, invalidRow, validRow2}

	result := MapAttendanceToManpower(rows)

	// Since row 2 is invalid, we should have 2 successful payloads and 1 failure
	assert.Len(t, result.Payloads, 2)
	assert.Len(t, result.Failures, 1)

	// Verify Failures map
	errMsg, exists := result.Failures["ATT-2"]
	assert.True(t, exists)
	assert.Contains(t, errMsg, "Missing mandatory field")

	// Verify payload mapping (Onsite vs Offsite)
	// Payload 0 corresponding to validRow1 (Onsite)
	payload0 := result.Payloads[0]
	assert.Equal(t, "ATT-1", payload0.InternalAttendanceID)
	assert.Equal(t, 1, *payload0.SubmissionEntity)
	assert.Equal(t, "REF-1", *payload0.ProjectReferenceNumber)
	assert.Nil(t, payload0.OffsiteFabricatorCompanyName) // Should be nil for Entity 1

	// Payload 1 corresponding to validRow2 (Offsite)
	payload1 := result.Payloads[1]
	assert.Equal(t, "ATT-3", payload1.InternalAttendanceID)
	assert.Equal(t, 2, *payload1.SubmissionEntity)
	assert.Equal(t, "Fabricator", *payload1.OffsiteFabricatorCompanyName)
	assert.Nil(t, payload1.ProjectReferenceNumber) // Should be nil for Entity 2
}
