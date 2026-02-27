package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sgbuildex/internal/api/middleware"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/apperrors"
	"sgbuildex/internal/pkg/idgen"
)

type SiteRepository struct {
	db *sql.DB
}

func NewSiteRepository(db *sql.DB) ports.SiteRepository {
	return &SiteRepository{db: db}
}

func (r *SiteRepository) Get(ctx context.Context, userID, id string) (*domain.Site, error) {
	query := `
		SELECT 
            s.site_id, s.user_id, s.site_name, s.location, 
            s.latitude, s.longitude, s.created_at, s.updated_at, u.user_name,
            (SELECT COUNT(*) FROM devices d WHERE d.site_id = s.site_id AND d.status != 'inactive') as device_count,
            (SELECT COUNT(*) FROM workers w WHERE w.current_project_id IN (SELECT project_id FROM projects p WHERE p.site_id = s.site_id)) as worker_count
		FROM sites s
		LEFT JOIN users u ON s.user_id = u.user_id`

	args := []interface{}{id}
	if middleware.IsVendor(ctx) {
		query += " WHERE s.site_id = ?"
	} else {
		query += " WHERE s.site_id = ? AND s.user_id = ?"
		args = append(args, userID)
	}

	var s domain.Site
	var scanUserID, loc, userName sql.NullString
	var lat, lng sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&s.ID, &scanUserID, &s.Name, &loc,
		&lat, &lng, &s.CreatedAt, &s.UpdatedAt, &userName,
		&s.DeviceCount, &s.WorkerCount,
	)
	if err == sql.ErrNoRows {
		return nil, apperrors.NewNotFound("site", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	if scanUserID.Valid {
		s.UserID = scanUserID.String
	}
	if loc.Valid {
		s.Location = loc.String
	}
	if userName.Valid {
		s.UserName = userName.String
	}
	if lat.Valid {
		s.Latitude = lat.Float64
	}
	if lng.Valid {
		s.Longitude = lng.Float64
	}

	return &s, nil
}

func (r *SiteRepository) List(ctx context.Context, userID string) ([]domain.Site, error) {
	query := `
        SELECT 
            s.site_id, s.user_id, s.site_name, s.location,
            s.latitude, s.longitude, s.created_at, s.updated_at, u.user_name,
            (SELECT COUNT(*) FROM devices d WHERE d.site_id = s.site_id AND d.status != 'inactive') as device_count,
            (SELECT COUNT(*) FROM workers w WHERE w.current_project_id IN (SELECT project_id FROM projects p WHERE p.site_id = s.site_id)) as worker_count
        FROM sites s
        LEFT JOIN users u ON s.user_id = u.user_id
        WHERE 1=1 `

	args := []interface{}{}
	log.Printf("[SECURITY] SiteRepository.List: userID='%s'", userID)

	if userID == "" && !middleware.IsVendor(ctx) {
		return nil, apperrors.NewPermissionDenied("user_id is required for multi-tenant isolation")
	}

	if !middleware.IsVendor(ctx) {
		query += " AND s.user_id = ?"
		args = append(args, userID)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []domain.Site
	for rows.Next() {
		var s domain.Site
		var uid, loc, userName sql.NullString
		var lat, lng sql.NullFloat64
		if err := rows.Scan(
			&s.ID, &uid, &s.Name, &loc,
			&lat, &lng, &s.CreatedAt, &s.UpdatedAt, &userName,
			&s.DeviceCount, &s.WorkerCount,
		); err != nil {
			return nil, err
		}
		if uid.Valid {
			s.UserID = uid.String
		}
		if loc.Valid {
			s.Location = loc.String
		}
		if userName.Valid {
			s.UserName = userName.String
		}
		if lat.Valid {
			s.Latitude = lat.Float64
		}
		if lng.Valid {
			s.Longitude = lng.Float64
		}

		sites = append(sites, s)
	}
	return sites, nil
}

func (r *SiteRepository) Create(ctx context.Context, s *domain.Site) error {
	id, err := idgen.GenerateNextID(r.db, "sites", "site_id", "site")
	if err != nil {
		return fmt.Errorf("failed to generate site ID: %w", err)
	}
	s.ID = id
	query := `INSERT INTO sites (
        site_id, user_id, site_name, location, latitude, longitude
    ) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = r.db.ExecContext(ctx, query,
		s.ID, s.UserID, s.Name, s.Location, s.Latitude, s.Longitude,
	)
	return err
}

func (r *SiteRepository) Update(ctx context.Context, s *domain.Site) error {
	query := `UPDATE sites SET 
        user_id=?, site_name=?, location=?, latitude=?, longitude=? 
        WHERE site_id=?`
	_, err := r.db.ExecContext(ctx, query,
		s.UserID, s.Name, s.Location, s.Latitude, s.Longitude, s.ID,
	)
	return err
}

func (r *SiteRepository) Delete(ctx context.Context, userID, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Unassign projects from this site
	if _, err := tx.ExecContext(ctx, "UPDATE projects SET site_id = NULL WHERE site_id = ? AND user_id = ?", id, userID); err != nil {
		return fmt.Errorf("failed to unassign projects: %w", err)
	}

	// 2. Unassign devices from this site
	if _, err := tx.ExecContext(ctx, "UPDATE devices SET site_id = NULL WHERE site_id = ? AND user_id = ?", id, userID); err != nil {
		return fmt.Errorf("failed to unassign devices: %w", err)
	}

	// 3. Delete the site
	res, err := tx.ExecContext(ctx, "DELETE FROM sites WHERE site_id = ? AND user_id = ?", id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return apperrors.NewNotFound("site", id)
	}

	return tx.Commit()
}
