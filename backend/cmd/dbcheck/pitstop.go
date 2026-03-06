package main

import (
	"database/sql"
	"fmt"
	"log"
	"sgbuildex/internal/pkg/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.LoadConfig()
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT pitstop_auth_id, regulator_name, regulator_id, on_behalf_of_id FROM pitstop_authorisations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Printf("%-20s | %-20s | %-15s | %-15s\n", "Auth ID", "Reg Name", "Reg ID", "OnBehalfOf")
	for rows.Next() {
		var id, rName string
		var rID, oID sql.NullString

		if err := rows.Scan(&id, &rName, &rID, &oID); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-20s | %-20s | %-15s | %-15s\n", id, rName, rID.String, oID.String)
	}
}
