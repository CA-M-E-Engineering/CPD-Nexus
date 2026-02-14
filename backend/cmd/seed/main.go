package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("--- SGBuildex Seeder Starting ---")

	err := godotenv.Load("../../.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Note: No .env file found, using system environment variables")
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
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Try multiple locations for seed file
	seedFile := "migrate/init_data.sql"
	if _, err := os.Stat(seedFile); os.IsNotExist(err) {
		seedFile = "../../migrate/init_data.sql"
	}
	content, err := ioutil.ReadFile(seedFile)
	if err != nil {
		log.Fatalf("Failed to read seed file %s: %v", seedFile, err)
	}

	log.Printf("Executing seed script: %s", seedFile)
	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatalf("Execution failed: %v", err)
	}

	log.Println("--- Seeding Completed Successfully! ---")
}
