package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type AnalyticsHandler struct {
	DB *sql.DB
}

func NewAnalyticsHandler(db *sql.DB) *AnalyticsHandler {
	return &AnalyticsHandler{DB: db}
}

func (h *AnalyticsHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "Missing tenant_id parameter", http.StatusBadRequest)
		return
	}
	stats := map[string]interface{}{
		"total_workers":   0,
		"active_sites":    0,
		"active_projects": 0,
		"total_devices":   0,
		"compliance_rate": 95, // Mock for now, but localized
	}

	var totalWorkers, activeSites, activeProjects, totalDevices int

	// Simple counts with tenant isolation
	h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role IN ('worker', 'pic', 'manager') AND tenant_id = ?", tenantID).Scan(&totalWorkers)
	h.DB.QueryRow("SELECT COUNT(*) FROM sites WHERE status='active' AND tenant_id = ?", tenantID).Scan(&activeSites)
	h.DB.QueryRow("SELECT COUNT(*) FROM projects WHERE status='active' AND tenant_id = ?", tenantID).Scan(&activeProjects)
	h.DB.QueryRow("SELECT COUNT(*) FROM devices WHERE tenant_id = ?", tenantID).Scan(&totalDevices)

	stats["total_workers"] = totalWorkers
	stats["active_sites"] = activeSites
	stats["active_projects"] = activeProjects
	stats["total_devices"] = totalDevices

	json.NewEncoder(w).Encode(stats)
}

func (h *AnalyticsHandler) GetActivityLog(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "Missing tenant_id parameter", http.StatusBadRequest)
		return
	}
	// Mock activity log for MVP but keep it tied to tenant context if we had a table
	logs := []map[string]interface{}{
		{"id": 1, "user": "System", "action": "Tenant Dashboard Loaded", "target": tenantID, "time": "Just now"},
		{"id": 2, "user": "System", "action": "Daily Sync Check", "target": "Cloud", "time": "5 mins ago"},
	}
	json.NewEncoder(w).Encode(logs)
}

func (h *AnalyticsHandler) GetDetailedAnalytics(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "Missing tenant_id parameter", http.StatusBadRequest)
		return
	}

	response := make(map[string]interface{})

	// 1. Worker Distribution by Trade
	tradeRows, err := h.DB.Query("SELECT trade_code, COUNT(*) FROM users WHERE role IN ('worker', 'pic') AND tenant_id = ? GROUP BY trade_code", tenantID)
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
	statusRows, err := h.DB.Query("SELECT status, COUNT(*) FROM users WHERE role IN ('worker', 'pic') AND tenant_id = ? GROUP BY status", tenantID)
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
	deviceRows, err := h.DB.Query("SELECT status, COUNT(*) FROM devices WHERE tenant_id = ? GROUP BY status", tenantID)
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
	// For real implementation: SELECT DATE(check_in_time), COUNT(*) FROM attendance ... GROUP BY DATE
	response["attendance_trends"] = map[string]int{
		"Mon": 12, "Tue": 15, "Wed": 18, "Thu": 14, "Fri": 16, "Sat": 8, "Sun": 2,
	}

	json.NewEncoder(w).Encode(response)
}
