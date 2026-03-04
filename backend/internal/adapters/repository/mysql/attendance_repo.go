package mysql

import (
	"context"
	"database/sql"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/apperrors"
)

const (
	attendanceSelectFields = `
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.user_id,
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			s.site_name, s.location,
			p.project_reference_number, p.project_title, p.project_location_description,
			p.project_contract_number, p.project_contract_name, p.hdb_precinct_name,
			p.main_contractor_name, p.main_contractor_uen,
			w.name AS worker_name, w.person_id_no, w.person_id_and_work_pass_type, w.person_nationality, w.person_trade AS worker_trade,
			p.worker_company_name, p.worker_company_uen, p.worker_company_trade,
			p.worker_company_client_name, p.worker_company_client_uen,
			pic.name AS pic_name, pic.person_id_no AS pic_fin,
			pa.regulator_id, pa.regulator_name, pa.on_behalf_of_id
	`
	attendanceJoinBlock = `
		FROM attendance a
		JOIN sites s ON a.site_id = s.site_id
		JOIN workers w ON a.worker_id = w.worker_id
		LEFT JOIN projects p ON w.current_project_id = p.project_id
		LEFT JOIN pitstop_authorisations pa ON p.pitstop_auth_id = pa.pitstop_auth_id
		LEFT JOIN workers pic ON p.project_id = pic.current_project_id AND pic.role = 'pic'
	`
)

type AttendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository(db *sql.DB) ports.AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) Get(ctx context.Context, userID, id string) (*domain.Attendance, error) {
	query := `
		SELECT 
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.user_id, 
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			w.name as worker_name, s.site_name, a.created_at, a.updated_at
		FROM attendance a
		LEFT JOIN workers w ON a.worker_id = w.worker_id
		LEFT JOIN sites s ON a.site_id = s.site_id
		WHERE a.attendance_id = ? AND a.user_id = ?`

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

func (r *AttendanceRepository) List(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error) {
	query := `
		SELECT 
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.user_id, 
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			w.name as worker_name, s.site_name, a.created_at, a.updated_at
		FROM attendance a
		LEFT JOIN workers w ON a.worker_id = w.worker_id
		LEFT JOIN sites s ON a.site_id = s.site_id`
	args := []interface{}{}

	query += " WHERE a.user_id = ?"
	args = append(args, userID)

	if userID == "" {
		return nil, apperrors.NewPermissionDenied("user_id is required for multi-tenant isolation")
	}
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
	return records, nil
}

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

func (r *AttendanceRepository) ExtractPendingAttendance(ctx context.Context) ([]domain.AttendanceRow, error) {
	query := `SELECT ` + attendanceSelectFields + attendanceJoinBlock + `
		WHERE a.status != 'submitted'
		ORDER BY a.submission_date, a.attendance_id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.AttendanceRow
	for rows.Next() {
		var res domain.AttendanceRow
		var mcName, mcUEN, wcName, wcUEN, wcTrade, wccName, wccUEN, picName, picFIN sql.NullString
		var pTitle, pLoc, pCNo, pCName, pHDB, wPassType, pNat, regID, regName, obID sql.NullString
		err := rows.Scan(
			&res.AttendanceID,
			&res.DeviceID,
			&res.WorkerID,
			&res.SiteID,
			&res.UserID,
			&res.TimeIn,
			&res.TimeOut,
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
			&picName,
			&picFIN,
			&regID,
			&regName,
			&obID,
		)
		if err != nil {
			return nil, err
		}
		if pTitle.Valid {
			res.ProjectTitle = pTitle.String
		}
		if pLoc.Valid {
			res.ProjectLocation = pLoc.String
		}
		if pCNo.Valid {
			res.ProjectContractNo = pCNo.String
		}
		if pCName.Valid {
			res.ProjectContractName = pCName.String
		}
		if pHDB.Valid {
			res.HDBPrecinctName = pHDB.String
		}
		if wPassType.Valid {
			res.WorkerWorkPassType = wPassType.String
		}
		if pNat.Valid {
			res.WorkerNationality = pNat.String
		}
		if regID.Valid {
			res.RegulatorID = regID.String
		}
		if regName.Valid {
			res.RegulatorName = regName.String
		}
		if obID.Valid {
			res.OnBehalfOfID = obID.String
		}
		if mcName.Valid {
			res.SiteOwnerName = mcName.String
		}
		if mcUEN.Valid {
			res.SiteOwnerUEN = mcUEN.String
		}
		if wcName.Valid {
			res.EmployerName = wcName.String
		}
		if wcUEN.Valid {
			res.EmployerUEN = wcUEN.String
		}
		if wcTrade.Valid {
			res.EmployerTrade = wcTrade.String
		}
		if wccName.Valid {
			res.EmployerClientName = wccName.String
		}
		if wccUEN.Valid {
			res.EmployerClientUEN = wccUEN.String
		}
		if picName.Valid {
			res.PICName = picName.String
		}
		if picFIN.Valid {
			res.PICFIN = picFIN.String
		}

		results = append(results, res)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *AttendanceRepository) ExtractPendingAttendanceByProject(ctx context.Context, projectID string) ([]domain.AttendanceRow, error) {
	query := `SELECT ` + attendanceSelectFields + attendanceJoinBlock + `
		WHERE a.status != 'submitted' AND w.current_project_id = ?
		ORDER BY a.submission_date, a.attendance_id
	`

	rows, err := r.db.QueryContext(ctx, query, projectID)
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

	return results, nil
}

// ExtractProjectsWithPendingAttendance returns a list of unique projects that have pending attendance records
func (r *AttendanceRepository) ExtractProjectsWithPendingAttendance(ctx context.Context) ([]domain.Project, error) {
	query := `
		SELECT DISTINCT
			p.project_id, p.site_id, p.user_id, p.project_title, p.status, 
			p.project_reference_number, p.project_contract_number, p.project_contract_name,
			p.project_location_description, p.hdb_precinct_name,
			p.pitstop_auth_id, pa.dataset_name as pitstop_auth_name,
			p.main_contractor_name, p.main_contractor_uen,
			p.worker_company_name, p.worker_company_uen,
			p.worker_company_client_name, p.worker_company_client_uen,
			p.worker_company_trade,
			s.site_name
		FROM attendance a
		JOIN workers w ON a.worker_id = w.worker_id
		JOIN projects p ON w.current_project_id = p.project_id
		JOIN sites s ON p.site_id = s.site_id
		LEFT JOIN pitstop_authorisations pa ON p.pitstop_auth_id = pa.pitstop_auth_id
		WHERE a.status != 'submitted'
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		var title, ref, status, cRef, cName, loc, hdb, pAuthID, pAuthName, siteName sql.NullString
		var mcName, mcUEN, wcName, wcUEN, wccName, wccUEN, wcTrade sql.NullString

		err := rows.Scan(
			&p.ID,
			&p.SiteID,
			&p.UserID,
			&title,
			&status,
			&ref,
			&cRef,
			&cName,
			&loc,
			&hdb,
			&pAuthID,
			&pAuthName,
			&mcName,
			&mcUEN,
			&wcName,
			&wcUEN,
			&wccName,
			&wccUEN,
			&wcTrade,
			&siteName,
		)
		if err != nil {
			return nil, err
		}

		if title.Valid {
			p.Title = title.String
		}
		if status.Valid {
			p.Status = status.String
		}
		if ref.Valid {
			p.Reference = ref.String
		}
		if cRef.Valid {
			p.ContractRef = cRef.String
		}
		if cName.Valid {
			p.ContractName = cName.String
		}
		if loc.Valid {
			p.Location = loc.String
		}
		if hdb.Valid {
			p.HDBPrecinct = hdb.String
		}
		if pAuthID.Valid {
			p.PitstopAuthID = &pAuthID.String
		}
		if pAuthName.Valid {
			p.PitstopAuthName = &pAuthName.String
		}
		if mcName.Valid {
			p.MainContractorName = mcName.String
		}
		if mcUEN.Valid {
			p.MainContractorUEN = mcUEN.String
		}
		if wcName.Valid {
			p.WorkerCompanyName = wcName.String
		}
		if wcUEN.Valid {
			p.WorkerCompanyUEN = wcUEN.String
		}
		if wccName.Valid {
			p.WorkerCompanyClientName = wccName.String
		}
		if wccUEN.Valid {
			p.WorkerCompanyClientUEN = wccUEN.String
		}
		if wcTrade.Valid {
			p.WorkerCompanyTrade = wcTrade.String
		}
		if siteName.Valid {
			p.SiteName = siteName.String
		}

		projects = append(projects, p)
	}

	return projects, nil
}

func (r *AttendanceRepository) scanAttendanceRow(rows *sql.Rows) (domain.AttendanceRow, error) {
	var res domain.AttendanceRow
	var timeOut sql.NullTime
	var mcName, mcUEN, wcName, wcUEN, wcTrade, wccName, wccUEN, picName, picFIN sql.NullString
	var pTitle, pLoc, pCNo, pCName, pHDB, wPassType, pNat, regID, regName, obID sql.NullString

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
		&picName,
		&picFIN,
		&regID,
		&regName,
		&obID,
	)
	if err != nil {
		return res, err
	}

	// Convert sql.NullTime → *time.Time so the domain stays free of sql types
	if timeOut.Valid {
		t := timeOut.Time
		res.TimeOut = &t
	}

	// Helper to set string only if valid
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
	mapNull(picName, &res.PICName)
	mapNull(picFIN, &res.PICFIN)
	mapNull(regID, &res.RegulatorID)
	mapNull(regName, &res.RegulatorName)
	mapNull(obID, &res.OnBehalfOfID)

	return res, nil
}
