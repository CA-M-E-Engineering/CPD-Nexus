package sgbuildex

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // or whichever driver you use
)

// AttendanceRow represents a row fetched from the DB with joined info
type AttendanceRow struct {
	AttendanceID string
	DeviceID     string
	WorkerID     string
	SiteID       string
	TenantID     string
	TimeIn       time.Time
	TimeOut      sql.NullTime
	Direction    string
	TradeCode    string
	Status       string

	// Joined Site fields
	SiteName     string
	SiteLocation string
	ProjectRef   string

	// Offsite Fabricator fields (from Project)
	OffsiteFabricator         sql.NullString
	OffsiteFabricatorUEN      sql.NullString
	OffsiteFabricatorLocation sql.NullString

	// Main Contractor fields (from Project)
	SiteOwnerName string
	SiteOwnerUEN  string

	// Joined Worker fields
	WorkerName   string
	WorkerFIN    string
	WorkerTrade  string
	EmployerName string
	EmployerUEN  string

	SubmissionDate time.Time
}

// ExtractPendingAttendance fetches all attendance rows with status 'pending'
// joining with sites, projects, and users to get full context for SGBuildex.
func ExtractPendingAttendance(ctx context.Context, db *sql.DB) ([]AttendanceRow, error) {
	query := `
		SELECT 
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.tenant_id,
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			s.site_name, s.location,
			p.project_reference_number,
			p.offsite_fabricator_name, p.offsite_fabricator_uen, p.offsite_fabricator_location,
			p.main_contractor_name, p.main_contractor_uen,
			u.name AS worker_name, u.fin_nric, u.trade_code AS worker_trade,
			p.worker_company_name, p.worker_company_uen
		FROM attendance a
		JOIN sites s ON a.site_id = s.site_id
		JOIN users u ON a.worker_id = u.user_id
		LEFT JOIN projects p ON u.current_project_id = p.project_id
		WHERE a.status = 'pending'
		ORDER BY a.submission_date, a.attendance_id
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []AttendanceRow
	for rows.Next() {
		var r AttendanceRow
		var mcName, mcUEN, wcName, wcUEN sql.NullString
		err := rows.Scan(
			&r.AttendanceID,
			&r.DeviceID,
			&r.WorkerID,
			&r.SiteID,
			&r.TenantID,
			&r.TimeIn,
			&r.TimeOut,
			&r.Direction,
			&r.TradeCode,
			&r.Status,
			&r.SubmissionDate,
			&r.SiteName,
			&r.SiteLocation,
			&r.ProjectRef,
			&r.OffsiteFabricator,
			&r.OffsiteFabricatorUEN,
			&r.OffsiteFabricatorLocation,
			&mcName,
			&mcUEN,
			&r.WorkerName,
			&r.WorkerFIN,
			&r.WorkerTrade,
			&wcName,
			&wcUEN,
		)
		if err != nil {
			return nil, err
		}
		if mcName.Valid {
			r.SiteOwnerName = mcName.String
		}
		if mcUEN.Valid {
			r.SiteOwnerUEN = mcUEN.String
		}
		if wcName.Valid {
			r.EmployerName = wcName.String
		}
		if wcUEN.Valid {
			r.EmployerUEN = wcUEN.String
		}

		results = append(results, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// MonthlyDistributionRow represents aggregated attendance data for MD submissions
type MonthlyDistributionRow struct {
	FabricatorName     string
	FabricatorUEN      string
	FabricatorLocation string
	ProjectRef         string
	ProjectTitle       string
	ProjectLocation    string
	SubmissionMonth    string
	AttendanceCount    int
}

// ExtractMonthlyDistributionData fetches aggregated attendance counts for the current month
// for all projects that have an offsite fabricator assigned.
func ExtractMonthlyDistributionData(ctx context.Context, db *sql.DB) ([]MonthlyDistributionRow, error) {
	query := `
		SELECT 
			p.offsite_fabricator_name, p.offsite_fabricator_uen, p.offsite_fabricator_location,
			p.project_reference_number, p.project_title, p.project_location_description,
			DATE_FORMAT(a.submission_date, '%Y-%m') as submission_month,
			COUNT(*) as attendance_count
		FROM attendance a
		JOIN users u ON a.worker_id = u.user_id
		JOIN projects p ON u.current_project_id = p.project_id
		WHERE p.offsite_fabricator_uen IS NOT NULL 
		  AND p.offsite_fabricator_uen != ''
		  AND DATE_FORMAT(a.submission_date, '%Y-%m') = DATE_FORMAT(CURRENT_DATE, '%Y-%m')
		GROUP BY 
			p.offsite_fabricator_name, p.offsite_fabricator_uen, p.offsite_fabricator_location,
			p.project_reference_number, p.project_title, p.project_location_description,
			submission_month
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MonthlyDistributionRow
	for rows.Next() {
		var r MonthlyDistributionRow
		var fabName, fabLoc, projRef, projTitle, projLoc sql.NullString
		err := rows.Scan(
			&fabName,
			&r.FabricatorUEN,
			&fabLoc,
			&projRef,
			&projTitle,
			&projLoc,
			&r.SubmissionMonth,
			&r.AttendanceCount,
		)
		if err != nil {
			return nil, err
		}
		r.FabricatorName = fabName.String
		r.FabricatorLocation = fabLoc.String
		r.ProjectRef = projRef.String
		r.ProjectTitle = projTitle.String
		r.ProjectLocation = projLoc.String
		results = append(results, r)
	}

	return results, nil
}
