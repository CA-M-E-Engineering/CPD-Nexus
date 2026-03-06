package validation

import (
	"testing"
)

func TestValidateSubmissionMonth(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"2023-01", true},
		{"2024-12", true},
		{"2024-13", false},
		{"24-01", false},
		{"2024-1", false},
		{"abcd-ef", false},
	}

	for _, tt := range tests {
		if res := ValidateSubmissionMonth(tt.input); res != tt.expected {
			t.Errorf("ValidateSubmissionMonth(%s) = %v; want %v", tt.input, res, tt.expected)
		}
	}
}

func TestValidateUEN(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"201461234E", true},
		{"T14LL1234K", true},
		{"12345678A", true},
		{"53312345A", true},
		{"INVALID123", false},
		{"12345678", false},
	}

	for _, tt := range tests {
		if res := ValidateUEN(tt.input); res != tt.expected {
			t.Errorf("ValidateUEN(%s) = %v; want %v", tt.input, res, tt.expected)
		}
	}
}

func TestValidateNRICFIN(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"S1234567A", true},
		{"T1234567B", true},
		{"F1234567C", true},
		{"G1234567D", true},
		{"M1234567E", true},
		{"S1234567", false},
		{"ABC123456", false},
	}

	for _, tt := range tests {
		if res := ValidateNRICFIN(tt.input); res != tt.expected {
			t.Errorf("ValidateNRICFIN(%s) = %v; want %v", tt.input, res, tt.expected)
		}
	}
}

func TestValidateWorkPassType(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"SP", true},
		{"WP", true},
		{"EP", true},
		{"LTVP", true},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		if res := ValidateWorkPassType(tt.input); res != tt.expected {
			t.Errorf("ValidateWorkPassType(%s) = %v; want %v", tt.input, res, tt.expected)
		}
	}
}

func TestValidateNRICWithPassType(t *testing.T) {
	tests := []struct {
		nric     string
		passType string
		expected bool
	}{
		{"S1234567A", "SP", true},
		{"T1234567B", "SB", true},
		{"F1234567C", "EP", true},
		{"G1234567D", "WP", true},
		{"M1234567E", "LTVP", true},
		{"F1234567C", "SP", false},
		{"S1234567A", "WP", false},
		{"", "SP", true},
	}

	for _, tt := range tests {
		if res := ValidateNRICWithPassType(tt.nric, tt.passType); res != tt.expected {
			t.Errorf("ValidateNRICWithPassType(%s, %s) = %v; want %v", tt.nric, tt.passType, res, tt.expected)
		}
	}
}
