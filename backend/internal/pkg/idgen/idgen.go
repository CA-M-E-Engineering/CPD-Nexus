package idgen

import (
	"database/sql"
	"time"
)

// GenerateNextID creates a centralized date-based unique ID.
// It uses only the first letter of the given string prefix.
// Example: prefix "worker" -> "w20260225135067"
func GenerateNextID(db *sql.DB, table, column, prefix string) (string, error) {
	if len(prefix) == 0 {
		prefix = "x"
	}
	firstLetter := prefix[:1]

	id := firstLetter + time.Now().Format("20060102150405")

	return id, nil
}
