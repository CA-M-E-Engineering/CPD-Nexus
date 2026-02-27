package mysql

import (
	"context"
	"database/sql"
	"sgbuildex/internal/core/domain"
)

type MySQLSettingsRepository struct {
	DB *sql.DB
}

func NewMySQLSettingsRepository(db *sql.DB) *MySQLSettingsRepository {
	return &MySQLSettingsRepository{DB: db}
}

func (r *MySQLSettingsRepository) GetSettings(ctx context.Context) (*domain.SystemSettings, error) {
	query := `
		SELECT id, attendance_sync_time, cpd_submission_time, 
		       max_payload_size_kb, max_workers_per_request, max_requests_per_minute, updated_at 
		FROM system_settings WHERE id = 1`

	var s domain.SystemSettings
	var updated sql.NullTime
	var cpdTime, syncInterval string

	err := r.DB.QueryRowContext(ctx, query).Scan(
		&s.ID,
		&syncInterval,
		&cpdTime,
		&s.MaxPayloadSizeKB,
		&s.MaxWorkersPerRequest,
		&s.MaxRequestsPerMinute,
		&updated,
	)
	if err != nil {
		return nil, err
	}

	s.CPDSubmissionTime = cpdTime
	s.AttendanceSyncTime = syncInterval
	if updated.Valid {
		s.UpdatedAt = updated.Time
	}

	return &s, nil
}

func (r *MySQLSettingsRepository) UpdateSettings(ctx context.Context, s domain.SystemSettings) error {
	query := `
		UPDATE system_settings 
		SET attendance_sync_time=?, cpd_submission_time=?,
		    max_payload_size_kb=?, max_workers_per_request=?, max_requests_per_minute=?
		WHERE id=1`
	_, err := r.DB.ExecContext(ctx, query,
		s.AttendanceSyncTime,
		s.CPDSubmissionTime,
		s.MaxPayloadSizeKB,
		s.MaxWorkersPerRequest,
		s.MaxRequestsPerMinute,
	)
	return err
}

func (r *MySQLSettingsRepository) GetDeviceStats(ctx context.Context) (int, int, error) {
	var total, unassigned int

	// Total active devices
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM devices WHERE status != 'inactive'").Scan(&total)
	if err != nil {
		return 0, 0, err
	}

	// Unassigned active devices (Assigned to Vendor)
	query := `
		SELECT COUNT(*) 
		FROM devices d
		JOIN users u ON d.user_id = u.user_id
		WHERE d.status != 'inactive' AND u.user_type = 'vendor'
	`
	err = r.DB.QueryRowContext(ctx, query).Scan(&unassigned)
	if err != nil {
		return 0, 0, err
	}

	return total, total - unassigned, nil
}
