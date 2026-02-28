package mysql

import (
	"context"
	"database/sql"
	"sgbuildex/internal/core/ports"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) ports.AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) GetDashboardStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"total_workers":   0,
		"active_sites":    0,
		"active_projects": 0,
		"total_devices":   0,
		"compliance_rate": 95,
	}

	var totalWorkers, activeSites, activeProjects, totalDevices int
	var err error

	if userID == "tenant-vendor-1" {
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM workers WHERE status = 'active' AND role IN ('worker', 'pic', 'manager')").Scan(&totalWorkers)
		if err != nil {
			return nil, err
		}
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sites").Scan(&activeSites)
		if err != nil {
			return nil, err
		}
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM projects WHERE status='active'").Scan(&activeProjects)
		if err != nil {
			return nil, err
		}
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM devices").Scan(&totalDevices)
		if err != nil {
			return nil, err
		}
	} else {
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM workers WHERE status = 'active' AND role IN ('worker', 'pic', 'manager') AND user_id = ?", userID).Scan(&totalWorkers)
		if err != nil {
			return nil, err
		}
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sites WHERE user_id = ?", userID).Scan(&activeSites)
		if err != nil {
			return nil, err
		}
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM projects WHERE status='active' AND user_id = ?", userID).Scan(&activeProjects)
		if err != nil {
			return nil, err
		}
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM devices WHERE user_id = ?", userID).Scan(&totalDevices)
		if err != nil {
			return nil, err
		}
	}

	stats["total_workers"] = totalWorkers
	stats["active_sites"] = activeSites
	stats["active_projects"] = activeProjects
	stats["total_devices"] = totalDevices

	return stats, nil
}

func (r *AnalyticsRepository) GetDetailedAnalytics(ctx context.Context, userID string) (map[string]interface{}, error) {
	response := make(map[string]interface{})

	// 1. Worker Distribution by Trade
	var tradeRows *sql.Rows
	var err error
	if userID == "tenant-vendor-1" {
		tradeRows, err = r.db.QueryContext(ctx, "SELECT person_trade, COUNT(*) FROM workers WHERE status = 'active' AND role IN ('worker', 'pic') GROUP BY person_trade")
	} else {
		tradeRows, err = r.db.QueryContext(ctx, "SELECT person_trade, COUNT(*) FROM workers WHERE status = 'active' AND role IN ('worker', 'pic') AND user_id = ? GROUP BY person_trade", userID)
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
	if userID == "tenant-vendor-1" {
		statusRows, err = r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM workers WHERE status = 'active' AND role IN ('worker', 'pic') GROUP BY status")
	} else {
		statusRows, err = r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM workers WHERE status = 'active' AND role IN ('worker', 'pic') AND user_id = ? GROUP BY status", userID)
	}
	if err == nil {
		defer statusRows.Close()
		statuses := make(map[string]int)
		for statusRows.Next() {
			var status string
			var count int
			if err := statusRows.Scan(&status, &count); err == nil {
				statuses[status] = count
			}
		}
		response["workers_by_status"] = statuses
	}

	// 3. Device Status Distribution
	var deviceRows *sql.Rows
	if userID == "tenant-vendor-1" {
		deviceRows, err = r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM devices GROUP BY status")
	} else {
		deviceRows, err = r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM devices WHERE user_id = ? GROUP BY status", userID)
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

	// 4. Mock Attendance Trends (Last 7 days)
	response["attendance_trends"] = map[string]int{
		"Mon": 12, "Tue": 15, "Wed": 18, "Thu": 14, "Fri": 16, "Sat": 8, "Sun": 2,
	}

	return response, nil
}
