package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/idgen"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) ports.DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Get(ctx context.Context, id string) (*domain.Device, error) {
	query := `
		SELECT 
			d.device_id, d.sn, d.model, d.status, 
			s.site_name, d.site_id,
			u.user_name, d.user_id,
			d.last_heartbeat, d.last_online_check, d.battery
		FROM devices d
		LEFT JOIN sites s ON d.site_id = s.site_id
		LEFT JOIN users u ON d.user_id = u.user_id
		WHERE d.device_id = ?`

	var d domain.Device
	var siteName, siteID, userName, userID sql.NullString
	var lastBeat, lastCheck sql.NullTime
	var status sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.SN, &d.Model, &status,
		&siteName, &siteID, &userName, &userID,
		&lastBeat, &lastCheck, &d.Battery,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	if status.Valid {
		d.Status = domain.DeviceStatus(status.String)
	}
	if siteName.Valid {
		d.SiteName = siteName.String
	}
	if siteID.Valid {
		sid := siteID.String
		d.SiteID = &sid
	}
	if userName.Valid {
		d.UserName = userName.String
	}
	if userID.Valid {
		d.UserID = userID.String
	}
	if lastBeat.Valid {
		d.LastHeartbeat = &lastBeat.Time
	}
	if lastCheck.Valid {
		d.LastOnlineCheck = &lastCheck.Time
	}

	return &d, nil
}

func (r *DeviceRepository) GetBySN(ctx context.Context, sn string) (*domain.Device, error) {
	query := `
		SELECT 
			d.device_id, d.sn, d.model, d.status, 
			s.site_name, d.site_id,
			u.user_name, d.user_id,
			d.last_heartbeat, d.last_online_check, d.battery
		FROM devices d
		LEFT JOIN sites s ON d.site_id = s.site_id
		LEFT JOIN users u ON d.user_id = u.user_id
		WHERE d.sn = ? LIMIT 1`

	var d domain.Device
	var siteName, siteID, userName, userID sql.NullString
	var lastBeat, lastCheck sql.NullTime
	var status sql.NullString

	err := r.db.QueryRowContext(ctx, query, sn).Scan(
		&d.ID, &d.SN, &d.Model, &status,
		&siteName, &siteID, &userName, &userID,
		&lastBeat, &lastCheck, &d.Battery,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get device by sn: %w", err)
	}

	if status.Valid {
		d.Status = domain.DeviceStatus(status.String)
	}
	if siteName.Valid {
		d.SiteName = siteName.String
	}
	if siteID.Valid {
		val := siteID.String
		d.SiteID = &val
	}
	if userName.Valid {
		d.UserName = userName.String
	}
	if userID.Valid {
		d.UserID = userID.String
	}
	if lastBeat.Valid {
		d.LastHeartbeat = &lastBeat.Time
	}
	if lastCheck.Valid {
		d.LastOnlineCheck = &lastCheck.Time
	}

	return &d, nil
}

func (r *DeviceRepository) List(ctx context.Context, userID string) ([]domain.Device, error) {
	query := `
		SELECT 
			d.device_id, d.sn, d.model, d.status, 
			s.site_name, d.site_id,
			u.user_name, d.user_id,
			d.last_heartbeat, d.last_online_check, d.battery
		FROM devices d
		LEFT JOIN sites s ON d.site_id = s.site_id
		LEFT JOIN users u ON d.user_id = u.user_id
		WHERE d.status != 'inactive'`

	args := []interface{}{}
	if userID != "" {
		query += " AND d.user_id = ?"
		args = append(args, userID)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	defer rows.Close()

	var devices []domain.Device
	for rows.Next() {
		var d domain.Device
		var siteName, siteID, uName, uid sql.NullString
		var lastBeat, lastCheck sql.NullTime
		var status sql.NullString

		if err := rows.Scan(
			&d.ID, &d.SN, &d.Model, &status,
			&siteName, &siteID, &uName, &uid,
			&lastBeat, &lastCheck, &d.Battery,
		); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		if status.Valid {
			d.Status = domain.DeviceStatus(status.String)
		}
		if siteName.Valid {
			d.SiteName = siteName.String
		}
		if siteID.Valid {
			val := siteID.String
			d.SiteID = &val
		}
		if uName.Valid {
			d.UserName = uName.String
		}
		if uid.Valid {
			d.UserID = uid.String
		}
		if lastBeat.Valid {
			d.LastHeartbeat = &lastBeat.Time
		}
		if lastCheck.Valid {
			d.LastOnlineCheck = &lastCheck.Time
		}

		devices = append(devices, d)
	}
	return devices, nil
}

func (r *DeviceRepository) Create(ctx context.Context, d *domain.Device) error {
	id, err := idgen.GenerateNextID(r.db, "devices", "device_id", "device")
	if err != nil {
		return fmt.Errorf("failed to generate device ID: %w", err)
	}
	d.ID = id
	query := `INSERT INTO devices (device_id, sn, user_id, site_id, model, status, last_heartbeat) VALUES (?, ?, ?, ?, ?, ?, NOW())`
	_, err = r.db.ExecContext(ctx, query, d.ID, d.SN, d.UserID, d.SiteID, d.Model, d.Status)
	return err
}

func (r *DeviceRepository) Update(ctx context.Context, d *domain.Device) error {
	query := "UPDATE devices SET sn=?, model=?, status=?"
	args := []interface{}{d.SN, d.Model, d.Status}
	if d.SiteID != nil {
		query += ", site_id=?"
		args = append(args, *d.SiteID)
	}
	if d.UserID != "" {
		query += ", user_id=?"
		args = append(args, d.UserID)
	}
	query += " WHERE device_id=?"
	args = append(args, d.ID)
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE devices SET status = 'inactive' WHERE device_id = ?", id)
	return err
}

func (r *DeviceRepository) AssignToUser(ctx context.Context, userID string, deviceIDs []string) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE devices SET user_id = ?, site_id = NULL, status = 'offline' WHERE device_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, id := range deviceIDs {
		if _, err := stmt.ExecContext(ctx, userID, id); err != nil {
			return err
		}
	}
	return nil
}
