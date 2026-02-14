package mysql

import (
	"context"
	"database/sql"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type AttendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository(db *sql.DB) ports.AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) Get(ctx context.Context, id string) (*domain.Attendance, error) {
	query := `
		SELECT 
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.tenant_id, 
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			u.name as worker_name, s.site_name, a.created_at, a.updated_at
		FROM attendance a
		LEFT JOIN users u ON a.worker_id = u.user_id
		LEFT JOIN sites s ON a.site_id = s.site_id
		WHERE a.attendance_id = ?`

	var a domain.Attendance
	var timeIn, timeOut sql.NullTime
	var subDate, wName, sName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.DeviceID, &a.WorkerID, &a.SiteID, &a.TenantID,
		&timeIn, &timeOut, &a.Direction, &a.TradeCode, &a.Status, &subDate,
		&wName, &sName, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
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

func (r *AttendanceRepository) List(ctx context.Context, tenantID, siteID, workerID, date string) ([]domain.Attendance, error) {
	query := `
		SELECT 
			a.attendance_id, a.device_id, a.worker_id, a.site_id, a.tenant_id, 
			a.time_in, a.time_out, a.direction, a.trade_code, a.status, a.submission_date,
			u.name as worker_name, s.site_name, a.created_at, a.updated_at
		FROM attendance a
		LEFT JOIN users u ON a.worker_id = u.user_id
		LEFT JOIN sites s ON a.site_id = s.site_id
		WHERE 1=1
	`
	args := []interface{}{}

	if tenantID != "" {
		query += " AND a.tenant_id = ?"
		args = append(args, tenantID)
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
			&a.ID, &a.DeviceID, &a.WorkerID, &a.SiteID, &a.TenantID,
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
			attendance_id, device_id, worker_id, site_id, tenant_id,
			time_in, time_out, direction, trade_code, status,
			submission_date, response_payload, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.DeviceID, a.WorkerID, a.SiteID, a.TenantID,
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
