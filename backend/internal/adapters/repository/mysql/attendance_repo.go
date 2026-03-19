package mysql

import (
	"context"
	"database/sql"
	"time"

	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/pkg/apperrors"
)

// SQL fragments shared across extraction queries to avoid repetition.
const (
	attendanceSelectFields = `
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.user_id,
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			s.site_name, s.location,
			p.project_reference_number, p.project_title, p.project_location_description,
			p.project_contract_number, p.project_contract_name, p.hdb_precinct_name,
			p.submission_entity, p.offsite_fabricator_name, p.offsite_fabricator_uen, p.offsite_fabricator_location,
			p.main_contractor_name, p.main_contractor_uen,
			w.name AS worker_name, w.person_id_no, w.person_id_and_work_pass_type, w.person_nationality, w.person_trade AS worker_trade,
			p.worker_company_name, p.worker_company_uen, p.worker_company_trade,
			p.worker_company_client_name, p.worker_company_client_uen,
			pa.regulator_id, pa.regulator_name, pa.on_behalf_of_id
	`
	attendanceJoinBlock = `
		FROM attendance a
		JOIN sites s ON a.site_id = s.site_id
		JOIN workers w ON a.worker_id = w.worker_id
		LEFT JOIN projects p ON w.current_project_id = p.project_id
		LEFT JOIN pitstop_authorisations pa ON p.pitstop_auth_id = pa.pitstop_auth_id
	`
)

type AttendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository(db *sql.DB) ports.AttendanceRepository {
	return &AttendanceRepository{db: db}
}

// Get retrieves a single attendance record by ID, scoped to the given user.
func (r *AttendanceRepository) Get(ctx context.Context, userID, id string) (*domain.Attendance, error) {
	query := `
		SELECT
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.user_id,
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			w.name AS worker_name, s.site_name, a.created_at, a.updated_at
		FROM attendance a
		LEFT JOIN workers w ON a.worker_id = w.worker_id
		LEFT JOIN sites s ON a.site_id = s.site_id
		WHERE a.attendance_id = ? AND a.user_id = ?
	`

	var a domain.Attendance
	var timeIn, timeOut sql.NullTime
	var subDate, wName, sName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&a.ID, &a.DeviceID, &a.WorkerID, &a.SiteID, &a.UserID,
		&timeIn, &timeOut, &a.Direction, &a.TradeCode, &a.Status, &subDate,
		&wName, &sName, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperrors.NewNotFound("attendance", id)
	}
	if err != nil {
		return nil, err
	}

	if timeIn.Valid {
		a.TimeIn = &timeIn.Time
	}
	if timeOut.Valid {
		a.TimeOut = &timeOut.Time
	}
	if subDate.Valid {
		a.SubmissionDate = subDate.String
	}
	if wName.Valid {
		a.WorkerName = wName.String
	}
	if sName.Valid {
		a.SiteName = sName.String
	}

	return &a, nil
}

// List retrieves attendance records filtered by optional siteID, workerID, and date, scoped to userID.
func (r *AttendanceRepository) List(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error) {
	if userID == "" {
		return nil, apperrors.NewPermissionDenied("user_id is required for multi-tenant isolation")
	}

	query := `
		SELECT
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.user_id,
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			w.name AS worker_name, s.site_name, a.created_at, a.updated_at
		FROM attendance a
		LEFT JOIN workers w ON a.worker_id = w.worker_id
		LEFT JOIN sites s ON a.site_id = s.site_id
		WHERE a.user_id = ?`

	args := []interface{}{userID}

	if siteID != "" {
		query += " AND a.site_id = ?"
		args = append(args, siteID)
	}
	if workerID != "" {
		query += " AND a.worker_id = ?"
		args = append(args, workerID)
	}
	if date != "" {
		query += " AND DATE(a.time_in) = ?"
		args = append(args, date)
	}
	query += " ORDER BY a.time_in DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.Attendance
	for rows.Next() {
		var a domain.Attendance
		var timeIn, timeOut sql.NullTime
		var subDate, wName, sName sql.NullString

		if err := rows.Scan(
			&a.ID, &a.DeviceID, &a.WorkerID, &a.SiteID, &a.UserID,
			&timeIn, &timeOut, &a.Direction, &a.TradeCode, &a.Status, &subDate,
			&wName, &sName, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if timeIn.Valid {
			a.TimeIn = &timeIn.Time
		}
		if timeOut.Valid {
			a.TimeOut = &timeOut.Time
		}
		if subDate.Valid {
			a.SubmissionDate = subDate.String
		}
		if wName.Valid {
			a.WorkerName = wName.String
		}
		if sName.Valid {
			a.SiteName = sName.String
		}

		records = append(records, a)
	}

	return records, rows.Err()
}

// Create inserts a new attendance record into the database.
func (r *AttendanceRepository) Create(ctx context.Context, a *domain.Attendance) error {
	query := `
		INSERT INTO attendance (
			attendance_id, device_id, worker_id, site_id, user_id,
			time_in, time_out, direction, trade_code, status,
			submission_date, response_payload, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.DeviceID, a.WorkerID, a.SiteID, a.UserID,
		a.TimeIn, a.TimeOut, a.Direction, a.TradeCode, a.Status,
		a.SubmissionDate, a.ResponsePayload,
	)
	return err
}

// Update modifies the TimeIn and TimeOut of an existing attendance record.
func (r *AttendanceRepository) Update(ctx context.Context, userID, id string, timeIn, timeOut *time.Time) error {
	query := `
		UPDATE attendance 
		SET time_in = ?, time_out = ?, updated_at = NOW() 
		WHERE attendance_id = ? AND user_id = ?
	`
	res, err := r.db.ExecContext(ctx, query, timeIn, timeOut, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return apperrors.NewNotFound("attendance", id)
	}
	return nil
}

// GetMaxID returns the highest attendance_id matching the given LIKE pattern.
func (r *AttendanceRepository) GetMaxID(ctx context.Context, pattern string) (string, error) {
	var maxID sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT MAX(attendance_id)
		FROM attendance
		WHERE attendance_id LIKE ?
	`, pattern).Scan(&maxID)
	if err != nil {
		return "", err
	}
	return maxID.String, nil
}

// ExtractPendingAttendance returns all non-submitted attendance rows with full joined data for SGBuildex submission.
func (r *AttendanceRepository) ExtractPendingAttendance(ctx context.Context) ([]domain.AttendanceRow, error) {
	query := `SELECT ` + attendanceSelectFields + attendanceJoinBlock + `
		WHERE a.status != 'submitted'
		ORDER BY a.submission_date, a.attendance_id
	`
	return r.queryAttendanceRows(ctx, query)
}

// ExtractPendingAttendanceByProject returns non-submitted attendance rows for a specific project.
// If userID is empty, it bypasses the user filter (for Admins).
func (r *AttendanceRepository) ExtractPendingAttendanceByProject(ctx context.Context, userID, projectID string) ([]domain.AttendanceRow, error) {
	query := `SELECT ` + attendanceSelectFields + attendanceJoinBlock + `
		WHERE a.status != 'submitted' AND w.current_project_id = ?`

	args := []interface{}{projectID}

	if userID != "" {
		query += ` AND a.user_id = ?`
		args = append(args, userID)
	}

	query += ` ORDER BY a.submission_date, a.attendance_id`

	return r.queryAttendanceRows(ctx, query, args...)
}

// ExtractProjectsWithPendingAttendance returns distinct projects that have attendance records not yet submitted.
// If userID is empty, it bypasses the user filter (for Admins).
func (r *AttendanceRepository) ExtractProjectsWithPendingAttendance(ctx context.Context, userID string) ([]domain.Project, error) {
	query := `
		SELECT DISTINCT
			p.project_id, p.site_id, p.user_id, p.project_title, p.status,
			p.project_reference_number, p.project_contract_number, p.project_contract_name,
			p.project_location_description, p.hdb_precinct_name,
			p.submission_entity, p.offsite_fabricator_name, p.offsite_fabricator_uen, p.offsite_fabricator_location,
			p.pitstop_auth_id, pa.dataset_name AS pitstop_auth_name,
			p.main_contractor_name, p.main_contractor_uen,
			p.worker_company_name, p.worker_company_uen,
			p.worker_company_client_name, p.worker_company_client_uen,
			p.worker_company_trade,
			s.site_name
		FROM projects p
		JOIN sites s ON p.site_id = s.site_id
		LEFT JOIN pitstop_authorisations pa ON p.pitstop_auth_id = pa.pitstop_auth_id
		JOIN workers w ON p.project_id = w.current_project_id
		JOIN attendance a ON w.worker_id = a.worker_id
		WHERE p.pitstop_auth_id IS NOT NULL AND p.pitstop_auth_id != ''
		AND p.status = ? AND a.status IN ('pending', 'failed')
	`

	args := []interface{}{domain.StatusActive}

	if userID != "" {
		query += ` AND p.user_id = ?`
		args = append(args, userID)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapNull := func(ns sql.NullString) string {
		if ns.Valid {
			return ns.String
		}
		return ""
	}
	mapNullPtr := func(ns sql.NullString) *string {
		if ns.Valid {
			s := ns.String
			return &s
		}
		return nil
	}

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		var title, ref, status, cRef, cName, loc, hdb, pAuthID, pAuthName, siteName sql.NullString
		var mcName, mcUEN, wcName, wcUEN, wccName, wccUEN, wcTrade sql.NullString
		var ofEnt sql.NullInt64 // SubmissionEntity is int
		var ofStrName, ofStrUEN, ofStrLoc sql.NullString

		if err := rows.Scan(
			&p.ID, &p.SiteID, &p.UserID,
			&title, &status, &ref, &cRef, &cName, &loc, &hdb,
			&ofEnt, &ofStrName, &ofStrUEN, &ofStrLoc,
			&pAuthID, &pAuthName,
			&mcName, &mcUEN,
			&wcName, &wcUEN, &wccName, &wccUEN, &wcTrade,
			&siteName,
		); err != nil {
			return nil, err
		}

		p.Title = mapNull(title)
		p.Status = mapNull(status)
		p.Reference = mapNull(ref)
		p.ContractRef = mapNull(cRef)
		p.ContractName = mapNull(cName)
		p.Location = mapNull(loc)
		p.HDBPrecinct = mapNull(hdb)
		if ofEnt.Valid {
			p.SubmissionEntity = int(ofEnt.Int64)
		}
		p.OffsiteFabricatorName = mapNull(ofStrName)
		p.OffsiteFabricatorUEN = mapNull(ofStrUEN)
		p.OffsiteFabricatorLocation = mapNull(ofStrLoc)
		p.PitstopAuthID = mapNullPtr(pAuthID)
		p.PitstopAuthName = mapNullPtr(pAuthName)
		p.MainContractorName = mapNull(mcName)
		p.MainContractorUEN = mapNull(mcUEN)
		p.WorkerCompanyName = mapNull(wcName)
		p.WorkerCompanyUEN = mapNull(wcUEN)
		p.WorkerCompanyClientName = mapNull(wccName)
		p.WorkerCompanyClientUEN = mapNull(wccUEN)
		p.WorkerCompanyTrade = mapNull(wcTrade)
		p.SiteName = mapNull(siteName)

		projects = append(projects, p)
	}

	return projects, rows.Err()
}

// queryAttendanceRows executes a query for AttendanceRow results and scans each row using the shared helper.
func (r *AttendanceRepository) queryAttendanceRows(ctx context.Context, query string, args ...interface{}) ([]domain.AttendanceRow, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.AttendanceRow
	for rows.Next() {
		res, err := r.scanAttendanceRow(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	return results, rows.Err()
}

// scanAttendanceRow scans a single database row into an AttendanceRow domain object.
// sql.NullTime is converted to *time.Time and sql.NullString fields are mapped via the mapNull helper.
func (r *AttendanceRepository) scanAttendanceRow(rows *sql.Rows) (domain.AttendanceRow, error) {
	var res domain.AttendanceRow
	var timeOut sql.NullTime
	var mcName, mcUEN, wcName, wcUEN, wcTrade, wccName, wccUEN sql.NullString
	var pTitle, pLoc, pCNo, pCName, pHDB, wPassType, pNat, regID, regName, obID sql.NullString
	var ofEnt sql.NullInt64
	var ofName, ofUEN, ofLoc sql.NullString

	err := rows.Scan(
		&res.AttendanceID,
		&res.DeviceID,
		&res.WorkerID,
		&res.SiteID,
		&res.UserID,
		&res.TimeIn,
		&timeOut,
		&res.Direction,
		&res.TradeCode,
		&res.Status,
		&res.SubmissionDate,
		&res.SiteName,
		&res.SiteLocation,
		&res.ProjectRef,
		&pTitle,
		&pLoc,
		&pCNo,
		&pCName,
		&pHDB,
		&ofEnt,
		&ofName,
		&ofUEN,
		&ofLoc,
		&mcName,
		&mcUEN,
		&res.WorkerName,
		&res.WorkerFIN,
		&wPassType,
		&pNat,
		&res.WorkerTrade,
		&wcName,
		&wcUEN,
		&wcTrade,
		&wccName,
		&wccUEN,
		&regID,
		&regName,
		&obID,
	)
	if err != nil {
		return res, err
	}

	// Convert sql.NullTime → *time.Time so the domain remains free of sql types
	if timeOut.Valid {
		t := timeOut.Time
		res.TimeOut = &t
	}

	mapNull := func(ns sql.NullString, target *string) {
		if ns.Valid {
			*target = ns.String
		}
	}

	mapNull(pTitle, &res.ProjectTitle)
	mapNull(pLoc, &res.ProjectLocation)
	mapNull(pCNo, &res.ProjectContractNo)
	mapNull(pCName, &res.ProjectContractName)
	mapNull(pHDB, &res.HDBPrecinctName)
	mapNull(mcName, &res.SiteOwnerName)
	mapNull(mcUEN, &res.SiteOwnerUEN)
	mapNull(wPassType, &res.WorkerWorkPassType)
	mapNull(pNat, &res.WorkerNationality)
	mapNull(wcName, &res.EmployerName)
	mapNull(wcUEN, &res.EmployerUEN)
	mapNull(wcTrade, &res.EmployerTrade)
	mapNull(wccName, &res.EmployerClientName)
	mapNull(wccUEN, &res.EmployerClientUEN)
	mapNull(regID, &res.RegulatorID)
	mapNull(regName, &res.RegulatorName)
	mapNull(obID, &res.OnBehalfOfID)

	if ofEnt.Valid {
		res.SubmissionEntity = int(ofEnt.Int64)
	}
	mapNull(ofName, &res.OffsiteFabricatorName)
	mapNull(ofUEN, &res.OffsiteFabricatorUEN)
	mapNull(ofLoc, &res.OffsiteFabricatorLocation)

	return res, nil
}
