package main // change to main so it can be executable

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"sgbuildex/internal/pkg/config"

	_ "github.com/go-sql-driver/mysql"
)

// RunMigrations reads all .sql files in folder and executes them in order
func RunMigrations(db *sql.DB, folder string) error {
	files, err := filepath.Glob(folder + "/*.sql")
	if err != nil {
		return err
	}

	sort.Strings(files) // ensure correct order

	for _, file := range files {
		log.Printf("Processing file: %s", file)
		sqlBytes, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("failed to exec %s: %w", file, err)
		}
		log.Printf("Migrated: %s\n", file)
	}

	return nil
}

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	// Run migrations
	if err := RunMigrations(db, "./"); err != nil { // ./ assuming you're in migrate folder
		log.Fatal("Migration failed:", err)
	}

	log.Println("All migrations ran successfully!")
}
