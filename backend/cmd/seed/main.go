package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"cpd-nexus/internal/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	logger.Infof("--- SGBuildex Seeder Starting ---")

	err := godotenv.Load("../../.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			logger.Infof("Note: No .env file found, using system environment variables")
		}
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "bas_user"
	}
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "new_password"
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1:3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "bas_mvp"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true", dbUser, dbPass, dbHost, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Try multiple locations for seed file
	seedFile := "migrate/init_data.sql"
	if _, err := os.Stat(seedFile); os.IsNotExist(err) {
		seedFile = "../../migrate/init_data.sql"
	}
	content, err := ioutil.ReadFile(seedFile)
	if err != nil {
		logger.Fatalf("Failed to read seed file %s: %v", seedFile, err)
	}

	logger.Infof("Executing seed script: %s", seedFile)
	_, err = db.Exec(string(content))
	if err != nil {
		logger.Fatalf("Execution failed: %v", err)
	}

	logger.Infof("--- Seeding Completed Successfully! ---")
}
