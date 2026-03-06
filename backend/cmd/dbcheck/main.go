package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

	rows, err := db.Query("SELECT p.project_id, p.pitstop_auth_id, p.status as p_status, a.status as a_status, p.user_id, count(*) FROM projects p LEFT JOIN workers w ON p.project_id = w.current_project_id LEFT JOIN attendance a ON w.worker_id = a.worker_id GROUP BY p.project_id, p.pitstop_auth_id, p.status, a.status, p.user_id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	f, _ := os.Create("db_out.txt")
	defer f.Close()

	fmt.Fprintf(f, "%-20s | %-12s | %-10s | %-10s | %-36s | %5s\n", "Proj", "Auth", "PStat", "AStat", "User", "Count")
	for rows.Next() {
		var pid string
		var authID sql.NullString
		var pstatus string
		var astatus sql.NullString
		var uid string
		var count int
		if err := rows.Scan(&pid, &authID, &pstatus, &astatus, &uid, &count); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(f, "%-20s | %-20s | %-10s | %-10s | %-36s | %5d\n", pid, authID.String, pstatus, astatus.String, uid, count)
	}

	fmt.Fprintf(f, "\n=== SUBMISSION LOGS ===\n")
	logs, err := db.Query("SELECT data_element_id, internal_id, status, error_message, created_at FROM submission_logs ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}
	defer logs.Close()

	for logs.Next() {
		var (
			elemID    string
			refID     string
			status    string
			errMsg    sql.NullString
			createdAt string
		)
		if err := logs.Scan(&elemID, &refID, &status, &errMsg, &createdAt); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(f, "[%s] %s (Ref: %s) -> %s\nError: %s\n", createdAt, elemID, refID, status, errMsg.String)
	}

	fmt.Fprintf(f, "\n=== DEMO PROJECT ATTENDANCE ===\n")
	attRows, err := db.Query("SELECT a.attendance_id, a.worker_id, a.status, a.error_message FROM attendance a LIMIT 10")
	if err == nil {
		defer attRows.Close()
		for attRows.Next() {
			var attID, wID, stat string
			var errMsg sql.NullString
			attRows.Scan(&attID, &wID, &stat, &errMsg)
			fmt.Fprintf(f, "Att ID: %s | Worker: %s | Stat: %s | Error: %s\n", attID, wID, stat, errMsg.String)
		}
	}
}
