package validation_test

import (
	"testing"

	"sgbuildex/internal/pkg/validation"
)

func TestValidateNRICFIN_Valid(t *testing.T) {
	validFINs := []string{"S1234567A", "T9876543B", "F1234567X", "G8765432Z"}
	for _, fin := range validFINs {
		if !validation.ValidateNRICFIN(fin) {
			t.Errorf("expected %q to be a valid NRIC/FIN", fin)
		}
	}
}

func TestValidateNRICFIN_Invalid(t *testing.T) {
	invalidFINs := []string{"", "1234567A", "X123456789", "S12345", "ABCDEFGH"}
	for _, fin := range invalidFINs {
		if validation.ValidateNRICFIN(fin) {
			t.Errorf("expected %q to be invalid NRIC/FIN", fin)
		}
	}
}

func TestValidateUEN_Valid(t *testing.T) {
	validUENs := []string{"197800111E", "200200111A", "53378491B"}
	for _, uen := range validUENs {
		if !validation.ValidateUEN(uen) {
			t.Errorf("expected %q to be a valid UEN", uen)
		}
	}
}

func TestValidateUEN_Invalid(t *testing.T) {
	invalidUENs := []string{"", "12345", "ABCDEFGHIJ", "1234567890X1"}
	for _, uen := range invalidUENs {
		if validation.ValidateUEN(uen) {
			t.Errorf("expected %q to be invalid UEN", uen)
		}
	}
}

func TestValidateWorkPassType_Valid(t *testing.T) {
	validTypes := []string{"SP", "SB", "EP", "SPASS", "WP", "ENTREPASS", "LTVP"}
	for _, wpt := range validTypes {
		if !validation.ValidateWorkPassType(wpt) {
			t.Errorf("expected %q to be a valid work pass type", wpt)
		}
	}
}

func TestValidateWorkPassType_Invalid(t *testing.T) {
	invalidTypes := []string{"", "UNKNOWN", "PR", "DEPENDENT"}
	for _, wpt := range invalidTypes {
		if validation.ValidateWorkPassType(wpt) {
			t.Errorf("expected %q to be an invalid work pass type", wpt)
		}
	}
}
