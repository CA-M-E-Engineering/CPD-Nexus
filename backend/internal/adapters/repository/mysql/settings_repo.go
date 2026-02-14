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
	query := `SELECT id, device_sync_interval, cpd_submission_time, response_size_limit, updated_at FROM system_settings WHERE id = 1`

	var s domain.SystemSettings
	// Scan time into string for simplicity, or handle sql.NullTime
	// Actually DB driver handles time.Time
	var updated sql.NullTime
	var cpdTime, syncInterval string

	err := r.DB.QueryRowContext(ctx, query).Scan(
		&s.ID,
		&syncInterval,
		&cpdTime,
		&s.ResponseSizeLimit,
		&updated,
	)
	if err != nil {
		return nil, err
	}

	s.CPDSubmissionTime = cpdTime
	s.DeviceSyncInterval = syncInterval
	if updated.Valid {
		s.UpdatedAt = updated.Time
	}

	return &s, nil
}

func (r *MySQLSettingsRepository) UpdateSettings(ctx context.Context, s domain.SystemSettings) error {
	query := `UPDATE system_settings SET device_sync_interval=?, cpd_submission_time=?, response_size_limit=? WHERE id=1`
	_, err := r.DB.ExecContext(ctx, query, s.DeviceSyncInterval, s.CPDSubmissionTime, s.ResponseSizeLimit)
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
