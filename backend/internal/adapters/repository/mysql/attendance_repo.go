package mysql

import (
	"context"
	"database/sql"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/apperrors"
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
	query := `
		SELECT 
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.user_id,
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			s.site_name, s.location,
			p.project_reference_number,
			p.offsite_fabricator_name, p.offsite_fabricator_uen, p.offsite_fabricator_location,
			p.main_contractor_name, p.main_contractor_uen,
			w.name AS worker_name, w.person_id_no, w.person_trade AS worker_trade,
			p.worker_company_name, p.worker_company_uen, p.worker_company_trade,
			pic.name AS pic_name, pic.person_id_no AS pic_fin
		FROM attendance a
		JOIN sites s ON a.site_id = s.site_id
		JOIN workers w ON a.worker_id = w.worker_id
		LEFT JOIN projects p ON w.current_project_id = p.project_id
		LEFT JOIN workers pic ON p.project_id = pic.current_project_id AND pic.role = 'pic'
		WHERE a.status = 'pending'
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
		var mcName, mcUEN, wcName, wcUEN, wcTrade, picName, picFIN sql.NullString
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
			&res.OffsiteFabricator,
			&res.OffsiteFabricatorUEN,
			&res.OffsiteFabricatorLocation,
			&mcName,
			&mcUEN,
			&res.WorkerName,
			&res.WorkerFIN,
			&res.WorkerTrade,
			&wcName,
			&wcUEN,
			&wcTrade,
			&picName,
			&picFIN,
		)
		if err != nil {
			return nil, err
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

func (r *AttendanceRepository) ExtractMonthlyDistributionData(ctx context.Context) ([]domain.MonthlyDistributionRow, error) {
	query := `
		SELECT 
			p.offsite_fabricator_name, p.offsite_fabricator_uen, p.offsite_fabricator_location,
			p.project_reference_number, p.project_title, p.project_location_description,
			DATE_FORMAT(a.submission_date, '%Y-%m') as submission_month,
			COUNT(*) as attendance_count
		FROM attendance a
		JOIN workers w ON a.worker_id = w.worker_id
		JOIN projects p ON w.current_project_id = p.project_id
		WHERE p.offsite_fabricator_uen IS NOT NULL 
		  AND p.offsite_fabricator_uen != ''
		  AND DATE_FORMAT(a.submission_date, '%Y-%m') = DATE_FORMAT(CURRENT_DATE, '%Y-%m')
		GROUP BY 
			p.offsite_fabricator_name, p.offsite_fabricator_uen, p.offsite_fabricator_location,
			p.project_reference_number, p.project_title, p.project_location_description,
			submission_month
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.MonthlyDistributionRow
	for rows.Next() {
		var res domain.MonthlyDistributionRow
		var fabName, fabLoc, projRef, projTitle, projLoc sql.NullString
		err := rows.Scan(
			&fabName,
			&res.FabricatorUEN,
			&fabLoc,
			&projRef,
			&projTitle,
			&projLoc,
			&res.SubmissionMonth,
			&res.AttendanceCount,
		)
		if err != nil {
			return nil, err
		}
		res.FabricatorName = fabName.String
		res.FabricatorLocation = fabLoc.String
		res.ProjectRef = projRef.String
		res.ProjectTitle = projTitle.String
		res.ProjectLocation = projLoc.String
		results = append(results, res)
	}

	return results, nil
}
