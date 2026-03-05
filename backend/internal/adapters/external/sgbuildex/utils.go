package sgbuildex

import (
	"strings"
	"time"
)

// Ptr returns a pointer to the trimmed string, or nil if the string is empty or "null" (case-insensitive).
func Ptr(s string) *string {
	trimmed := strings.TrimSpace(s)
	lower := strings.ToLower(trimmed)
	if trimmed == "" || lower == "null" {
		return nil
	}
	return &trimmed
}

// PtrInt returns a pointer to the given int value.
func PtrInt(v int) *int {
	return &v
}

// ptrIntOrDefault returns a pointer to v if v is non-zero, otherwise a pointer to def.
func ptrIntOrDefault(v int, def int) *int {
	if v == 0 {
		return &def
	}
	return &v
}

// FormatOptionalTime formats a nullable time pointer to RFC3339, returning nil if input is nil.
func FormatOptionalTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
