package mysql

import (
	"context"
	"database/sql"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"strings"
	"time"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) ports.AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) GetDashboardStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	isVendor := ports.IsVendor(ctx)
	contextUserID := ports.GetUserID(ctx)

	// Security/Scope logic
	queryUserID := userID
	if !isVendor {
		queryUserID = contextUserID
	}

	response := map[string]interface{}{
		"total_workers":   0,
		"active_sites":    0,
		"active_projects": 0,
		"total_devices":   0,
		"compliance_rate": 95,
	}

	var totalWorkers, activeSites, activeProjects, totalDevices int
	var err error

	if queryUserID == "all" || queryUserID == "tenant-vendor-1" {
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM workers WHERE status = ?", domain.StatusActive).Scan(&totalWorkers)
		if err != nil { return nil, err }
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sites").Scan(&activeSites)
		if err != nil { return nil, err }
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM projects WHERE status= ?", domain.StatusActive).Scan(&activeProjects)
		if err != nil { return nil, err }
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM devices").Scan(&totalDevices)
		if err != nil { return nil, err }
	} else {
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM workers WHERE status = ? AND user_id = ?", domain.StatusActive, queryUserID).Scan(&totalWorkers)
		if err != nil { return nil, err }
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sites WHERE user_id = ?", queryUserID).Scan(&activeSites)
		if err != nil { return nil, err }
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM projects WHERE status= ? AND user_id = ?", domain.StatusActive, queryUserID).Scan(&activeProjects)
		if err != nil { return nil, err }
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM devices WHERE user_id = ?", queryUserID).Scan(&totalDevices)
		if err != nil { return nil, err }
	}

	response["total_workers"] = totalWorkers
	response["active_sites"] = activeSites
	response["active_projects"] = activeProjects
	response["total_devices"] = totalDevices

	return response, nil
}

func (r *AnalyticsRepository) GetDetailedAnalytics(ctx context.Context, userID string) (map[string]interface{}, error) {
	isVendor := ports.IsVendor(ctx)
	contextUserID := ports.GetUserID(ctx)

	queryUserID := userID
	if !isVendor {
		queryUserID = contextUserID
	}

	response := make(map[string]interface{})

	// 1. Worker Distribution by Trade
	var tradeRows *sql.Rows
	var err error
	if queryUserID == "all" || queryUserID == "tenant-vendor-1" {
		tradeRows, err = r.db.QueryContext(ctx, "SELECT person_trade, COUNT(*) FROM workers WHERE status = ? GROUP BY person_trade", domain.StatusActive)
	} else {
		tradeRows, err = r.db.QueryContext(ctx, "SELECT person_trade, COUNT(*) FROM workers WHERE status = ? AND user_id = ? GROUP BY person_trade", domain.StatusActive, queryUserID)
	}
	if err == nil {
		defer tradeRows.Close()
		trades := make(map[string]int)
		for tradeRows.Next() {
			var code sql.NullString
			var count int
			if err := tradeRows.Scan(&code, &count); err == nil {
				label := "General"
				if code.Valid && code.String != "" {
					label = code.String
				}
				trades[label] = count
			}
		}
		response["workers_by_trade"] = trades
	}

	// 2. Worker Status Distribution
	var statusRows *sql.Rows
	if queryUserID == "all" || queryUserID == "tenant-vendor-1" {
		statusRows, err = r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM workers GROUP BY status")
	} else {
		statusRows, err = r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM workers WHERE user_id = ? GROUP BY status", queryUserID)
	}
	if err == nil {
		defer statusRows.Close()
		statuses := make(map[string]int)
		for statusRows.Next() {
			var st string
			var count int
			if err := statusRows.Scan(&st, &count); err == nil {
				statuses[st] = count
			}
		}
		response["workers_by_status"] = statuses
	}

	// 3. Device Status Distribution
	var deviceRows *sql.Rows
	if queryUserID == "all" || queryUserID == "tenant-vendor-1" {
		deviceRows, err = r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM devices GROUP BY status")
	} else {
		deviceRows, err = r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM devices WHERE user_id = ? GROUP BY status", queryUserID)
	}
	if err == nil {
		defer deviceRows.Close()
		dStatuses := make(map[string]int)
		for deviceRows.Next() {
			var status string
			var count int
			if err := deviceRows.Scan(&status, &count); err == nil {
				dStatuses[status] = count
			}
		}
		response["devices_by_status"] = dStatuses
	}

	// 4. Attendance Trends (Last 7 days)
	trends := make(map[string]int)
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	
	trendQuery := `
		SELECT DAYNAME(created_at) as day_name, COUNT(*) 
		FROM attendance 
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)`
	
	var trendArgs []interface{}
	if queryUserID != "all" && queryUserID != "tenant-vendor-1" && queryUserID != "" {
		trendQuery += " AND user_id = ?"
		trendArgs = append(trendArgs, queryUserID)
	}
	trendQuery += " GROUP BY day_name"

	trendRows, err := r.db.QueryContext(ctx, trendQuery, trendArgs...)
	if err == nil {
		defer trendRows.Close()
		for trendRows.Next() {
			var day string
			var count int
			if err := trendRows.Scan(&day, &count); err == nil {
				shortDay := day[:3]
				trends[shortDay] = count
			}
		}
	}
	
	finalTrends := make(map[string]int)
	for _, d := range days {
		finalTrends[d] = trends[d]
	}
	response["attendance_trends"] = finalTrends

	return response, nil
}

func (r *AnalyticsRepository) GetActivityLog(ctx context.Context, userID string, filters map[string]interface{}) ([]map[string]interface{}, error) {
	var rows *sql.Rows
	var err error

	isVendor := ports.IsVendor(ctx)
	contextUserID := ports.GetUserID(ctx)

	queryUserID := userID
	if !isVendor {
		queryUserID = contextUserID
	}

	query := `SELECT id, user_id, user_name, action, target_type, target_id, details, created_at 
			 FROM activity_logs WHERE 1=1`
	var args []interface{}

	if !isVendor || (queryUserID != "" && queryUserID != "all" && queryUserID != "tenant-vendor-1") {
		query += ` AND user_id = ?`
		args = append(args, queryUserID)
	}

	if action, ok := filters["action"].(string); ok && action != "" {
		query += ` AND action = ?`
		args = append(args, action)
	}

	if targetType, ok := filters["target_type"].(string); ok && targetType != "" {
		query += ` AND target_type = ?`
		args = append(args, targetType)
	}

	query += ` ORDER BY created_at DESC LIMIT 100`

	rows, err = r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var uid, uname, action, tType, tID, details sql.NullString
		var createdAt sql.NullTime

		if err := rows.Scan(&id, &uid, &uname, &action, &tType, &tID, &details, &createdAt); err != nil {
			return nil, err
		}

		log := map[string]interface{}{
			"id":          id,
			"user_id":     uid.String,
			"user_name":   uname.String,
			"action":      action.String,
			"target_type": tType.String,
			"target_id":   tID.String,
			"details":     details.String,
			"time":        "Just now",
		}

		if createdAt.Valid {
			log["created_at"] = createdAt.Time
			log["time"] = formatTimeAgo(createdAt.Time)
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *AnalyticsRepository) LogActivity(ctx context.Context, log map[string]interface{}) error {
	query := `INSERT INTO activity_logs (user_id, user_name, action, target_type, target_id, details, ip_address) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		log["user_id"],
		log["user_name"],
		log["action"],
		log["target_type"],
		log["target_id"],
		log["details"],
		log["ip_address"],
	)
	return err
}

func formatTimeAgo(t interface{}) string {
	var ts time.Time
	switch v := t.(type) {
	case time.Time:
		ts = v
	case *time.Time:
		if v == nil { return "" }
		ts = *v
	default:
		return "Unknown"
	}

	diff := time.Since(ts)
	if diff.Seconds() < 60 {
		return "Just now"
	}
	if diff.Minutes() < 60 {
		return strings.Split(diff.Truncate(time.Minute).String(), "m")[0] + "m ago"
	}
	if diff.Hours() < 24 {
		return strings.Split(diff.Truncate(time.Hour).String(), "h")[0] + "h ago"
	}
	days := int(diff.Hours() / 24)
	if days < 30 {
		if days == 1 { return "Yesterday" }
		return strings.Split(diff.Truncate(24*time.Hour).String(), "h")[0] + "d ago"
	}
	return ts.Format("2006-01-02")
}
