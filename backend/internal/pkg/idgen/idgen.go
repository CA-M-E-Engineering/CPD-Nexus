package idgen

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

// GenerateNextID finds the next sequential ID for a given table, column, and prefix.
// It assumes IDs are in the format "prefix_XXX" (e.g. user_001).
// It returns "prefix_001" if no matching IDs are found.
func GenerateNextID(db *sql.DB, table, column, prefix string) (string, error) {
	// 1. Get all IDs that match the prefix
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIKE '%s_%%'", column, table, column, prefix)
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	maxID := -1
	regex := regexp.MustCompile(fmt.Sprintf(`^%s_(\d+)$`, prefix))

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}

		matches := regex.FindStringSubmatch(id)
		if len(matches) == 2 {
			num, err := strconv.Atoi(matches[1])
			if err == nil && num > maxID {
				maxID = num
			}
		}
	}

	nextID := maxID + 1

	// Final safety check: if the proposed ID somehow exists (e.g. out-of-band edit),
	// increment until we find a free one.
	for {
		candidate := fmt.Sprintf("%s_%03d", prefix, nextID)
		var exists int
		checkQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, column)
		err := db.QueryRow(checkQuery, candidate).Scan(&exists)
		if err != nil {
			return "", err
		}
		if exists == 0 {
			return candidate, nil
		}
		nextID++
	}
}
