package validation

import (
	"regexp"
	"strings"
)

var (
	reSubmissionMonth = regexp.MustCompile(`^(?:19[5-9][0-9]|20[0-9]{2}|2100)-(?:0[1-9]|1[0-2])$`)
	reProjectRef      = regexp.MustCompile(`^[AE]\d{4}-[A-Za-z0-9]{5}-\d{4}$`)
	reUEN             = regexp.MustCompile(`^(?:\d{8}[A-Z]|\d{4}\d{5}[A-Z]|[TSR]\d{2}[A-Z0-9]{2}\d{4}[A-Z])$`)
	reNRICFIN         = regexp.MustCompile(`^[STFGM]\d{7}[A-Z0-9]$`)
	reHDBContract     = regexp.MustCompile(`^D/\d{1,5}/\d{2}$`)
	reLTAContract     = regexp.MustCompile(`^[A-Za-z0-9 /()-]{1,20}$`)
)

func ValidateSubmissionEntity(v int) bool {
	return v == 1 || v == 2
}

func ValidateSubmissionMonth(v string) bool {
	return reSubmissionMonth.MatchString(v)
}

func ValidateProjectReferenceNumber(v string) bool {
	return reProjectRef.MatchString(v)
}

func ValidateUEN(v string) bool {
	return reUEN.MatchString(v)
}

func ValidateNRICFIN(v string) bool {
	return reNRICFIN.MatchString(v)
}

func ValidateWorkPassType(v string) bool {
	enums := []string{"SP", "SB", "EP", "SPASS", "WP", "ENTREPASS", "LTVP"}
	for _, e := range enums {
		if v == e {
			return true
		}
	}
	return false
}

func ValidatePersonTrade(v string) bool {
	// Simple check for XX.XX format and ranges
	// Allowable Range: 1.1-1.5; 2.1-2.8; 3.1-3.11, 4.1-4.6
	validTrades := map[string]bool{
		"1.1": true, "1.2": true, "1.3": true, "1.4": true, "1.5": true,
		"2.1": true, "2.2": true, "2.3": true, "2.4": true, "2.5": true, "2.6": true, "2.7": true, "2.8": true,
		"3.1": true, "3.2": true, "3.3": true, "3.4": true, "3.5": true, "3.6": true, "3.7": true, "3.8": true, "3.9": true, "3.10": true, "3.11": true,
		"4.1": true, "4.2": true, "4.3": true, "4.4": true, "4.5": true, "4.6": true,
	}
	return validTrades[v]
}

func ValidateHDBContractNumber(v string) bool {
	return len(v) <= 10 && reHDBContract.MatchString(v)
}

func ValidateLTAContractNumber(v string) bool {
	return len(v) <= 20 && reLTAContract.MatchString(v)
}

func SanitizeUEN(v string) string {
	return strings.ToUpper(strings.TrimSpace(v))
}
