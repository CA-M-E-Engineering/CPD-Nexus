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

	// Offsite Fabricator fields (from Site)
	OffsiteFabricator         sql.NullString
	OffsiteFabricatorUEN      sql.NullString
	OffsiteFabricatorLocation sql.NullString

	// Joined Tenant fields (Main Contractor)
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
// joining with sites, tenants, and users to get full context for SGBuildex.
func ExtractPendingAttendance(ctx context.Context, db *sql.DB) ([]AttendanceRow, error) {
	query := `
		SELECT 
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.tenant_id,
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			s.site_name, s.location, s.project_ref,
			s.offsite_fabricator, s.offsite_fabricator_uen, s.offsite_fabricator_location,
			sc.company_name AS site_owner_name, sc.uen AS site_owner_uen,
			u.name AS worker_name, u.fin_nric, u.trade_code AS worker_trade,
			ec.company_name AS employer_name, ec.uen AS employer_uen
		FROM attendance a
		JOIN sites s ON a.site_id = s.site_id
		LEFT JOIN companies sc ON s.tenant_id = sc.tenant_id AND sc.company_type = 'contractor'
		JOIN users u ON a.worker_id = u.user_id
		LEFT JOIN companies ec ON u.tenant_id = ec.tenant_id AND ec.company_type = 'contractor'
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
			&r.SiteOwnerName,
			&r.SiteOwnerUEN,
			&r.WorkerName,
			&r.WorkerFIN,
			&r.WorkerTrade,
			&r.EmployerName,
			&r.EmployerUEN,
		)
		if err != nil {
			return nil, err
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
// for all sites that have an offsite fabricator assigned.
func ExtractMonthlyDistributionData(ctx context.Context, db *sql.DB) ([]MonthlyDistributionRow, error) {
	query := `
		SELECT 
			s.offsite_fabricator, s.offsite_fabricator_uen, s.offsite_fabricator_location,
			s.project_ref, s.site_name, s.location,
			DATE_FORMAT(a.submission_date, '%Y-%m') as submission_month,
			COUNT(*) as attendance_count
		FROM attendance a
		JOIN sites s ON a.site_id = s.site_id
		WHERE s.offsite_fabricator_uen IS NOT NULL 
		  AND s.offsite_fabricator_uen != ''
		  AND DATE_FORMAT(a.submission_date, '%Y-%m') = DATE_FORMAT(CURRENT_DATE, '%Y-%m')
		GROUP BY 
			s.offsite_fabricator, s.offsite_fabricator_uen, s.offsite_fabricator_location,
			s.project_ref, s.site_name, s.location,
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
