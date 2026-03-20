package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"cpd-nexus/internal/pkg/config"
)

func main() {
	cfg := config.LoadConfig()
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	fmt.Println("Updating users table schema...")
	
	// Add columns individually to handle cases where some might already exist
	columns := map[string]string{
		"bridge_ws_url":     "VARCHAR(255) DEFAULT NULL",
		"bridge_auth_token": "VARCHAR(255) DEFAULT NULL",
		"bridge_status":    "VARCHAR(50) DEFAULT 'inactive'",
	}

	for col, spec := range columns {
		query := fmt.Sprintf("ALTER TABLE users ADD COLUMN %s %s", col, spec)
		fmt.Printf("Executing: %s\n", query)
		if _, err := db.Exec(query); err != nil {
			fmt.Printf("Skipping %s (likely already exists or error: %v)\n", col, err)
		} else {
			fmt.Printf("Column %s added successfully.\n", col)
		}
	}

	fmt.Println("Migration complete.")
}
