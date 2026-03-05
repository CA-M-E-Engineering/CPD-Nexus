package validation

import (
	"regexp"
	"strings"
)

// Pre-compiled regular expressions for SGBuildex field validation.
var (
	reSubmissionMonth = regexp.MustCompile(`^(?:19[5-9]\d|20\d{2}|2100)-(?:0[1-9]|1[0-2])$`)
	reProjectRef      = regexp.MustCompile(`^[AE]\d{4}-[A-Za-z0-9]{5}-\d{4}$`)
	reUEN             = regexp.MustCompile(`^(?:\d{8}[A-Z]|\d{9}[A-Z]|[TSR]\d{2}[A-Z0-9]{2}\d{4}[A-Z])$`)
	reNRICFIN         = regexp.MustCompile(`^[STFGM]\d{7}[A-Z0-9]$`)
	reHDBContract     = regexp.MustCompile(`^D/\d{1,5}/\d{2}$`)
	reLTAContract     = regexp.MustCompile(`^[A-Za-z0-9 /()-]{1,20}$`)
)

// validWorkPassTypes lists all accepted work pass type values.
var validWorkPassTypes = map[string]bool{
	"SP": true, "SB": true, "EP": true,
	"SPASS": true, "WP": true, "ENTREPASS": true, "LTVP": true,
}

// validPersonTrades lists all accepted BCA trade codes.
var validPersonTrades = map[string]bool{
	"1.1": true, "1.2": true, "1.3": true, "1.4": true, "1.5": true,
	"2.1": true, "2.2": true, "2.3": true, "2.4": true, "2.5": true,
	"2.6": true, "2.7": true, "2.8": true,
	"3.1": true, "3.2": true, "3.3": true, "3.4": true, "3.5": true,
	"3.6": true, "3.7": true, "3.8": true, "3.9": true, "3.10": true, "3.11": true,
	"4.1": true, "4.2": true, "4.3": true, "4.4": true, "4.5": true, "4.6": true,
}

// ValidateSubmissionEntity checks that the entity is either 1 (Onsite Builder) or 2 (Offsite Fabricator).
func ValidateSubmissionEntity(v int) bool {
	return v == 1 || v == 2
}

// ValidateSubmissionMonth checks for the YYYY-MM format required by the API.
func ValidateSubmissionMonth(v string) bool {
	return reSubmissionMonth.MatchString(v)
}

// ValidateProjectReferenceNumber validates project reference numbers (e.g. A1234-AB123-2022).
func ValidateProjectReferenceNumber(v string) bool {
	return reProjectRef.MatchString(v)
}

// ValidateUEN validates a Unique Entity Number (Singapore).
func ValidateUEN(v string) bool {
	return reUEN.MatchString(v)
}

// ValidateNRICFIN validates a Singapore NRIC or FIN number.
func ValidateNRICFIN(v string) bool {
	return reNRICFIN.MatchString(v)
}

// ValidateWorkPassType checks that the work pass type is one of the accepted enum values.
func ValidateWorkPassType(v string) bool {
	return validWorkPassTypes[v]
}

// ValidatePersonTrade checks that the trade code is a valid BCA trade.
func ValidatePersonTrade(v string) bool {
	return validPersonTrades[v]
}

// localPassTypes are pass types that require a Singapore citizen prefix (S or T).
var localPassTypes = map[string]bool{"SP": true, "SB": true}

// foreignPassTypes are pass types that require a foreign identification prefix (F, G, or M).
var foreignPassTypes = map[string]bool{
	"EP": true, "SPASS": true, "WP": true, "ENTREPASS": true, "LTVP": true,
}

// ValidateNRICWithPassType cross-checks the NRIC/FIN century prefix against the pass type.
// Per ICA/MOM spec:
//   - SP (Singapore Pink IC) / SB (Singapore Blue IC): prefix must be S or T
//   - EP / SPASS / WP / ENTREPASS / LTVP: prefix must be F, G, or M
//
// Returns true if the combination is valid or if either field is empty (validated separately).
func ValidateNRICWithPassType(nric, passType string) bool {
	nric = strings.ToUpper(strings.TrimSpace(nric))
	passType = strings.ToUpper(strings.TrimSpace(passType))
	if nric == "" || passType == "" {
		return true // Missing values are caught by individual field validators
	}
	prefix := rune(nric[0])
	if localPassTypes[passType] {
		return prefix == 'S' || prefix == 'T'
	}
	if foreignPassTypes[passType] {
		return prefix == 'F' || prefix == 'G' || prefix == 'M'
	}
	return true // Unknown pass type — don't fail here, let ValidateWorkPassType catch it
}

// ValidateHDBContractNumber validates HDB contract number format (D/NNNNN/YY).
func ValidateHDBContractNumber(v string) bool {
	return len(v) <= 10 && reHDBContract.MatchString(v)
}

// ValidateLTAContractNumber validates LTA contract number format (max 20 alphanumeric chars).
func ValidateLTAContractNumber(v string) bool {
	return len(v) <= 20 && reLTAContract.MatchString(v)
}

// SanitizeUEN normalises a UEN by trimming whitespace and uppercasing.
func SanitizeUEN(v string) string {
	return strings.ToUpper(strings.TrimSpace(v))
}
