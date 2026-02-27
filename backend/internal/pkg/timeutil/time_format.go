package timeutil

import (
	"strings"
	"time"
)

// CleanDateTime removes 'T' and 'Z' characters from a datetime string for MySQL compatibility
func CleanDateTime(ds string) string {
	if len(ds) >= 19 {
		tmp := strings.Replace(ds, "T", " ", 1)
		tmp = strings.Replace(tmp, "Z", "", 1)
		return tmp
	}
	return ds
}

// ToRFC3339 attempts to parse a datetime string in various formats and returns it in RFC3339
func ToRFC3339(ds string) string {
	if ds == "" {
		return ""
	}

	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, ds); err == nil {
			return parsed.Format(time.RFC3339)
		}
	}

	return ds
}
