package sgbuildex

import "time"

// Ptr returns a pointer to the given string, or nil if the string is empty.
func Ptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// FormatOptionalTime formats a nullable time pointer to RFC3339, returning empty string if nil.
func FormatOptionalTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
