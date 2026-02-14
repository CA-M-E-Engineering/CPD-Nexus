package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type SiteRepository struct {
	db *sql.DB
}

func NewSiteRepository(db *sql.DB) ports.SiteRepository {
	return &SiteRepository{db: db}
}

func (r *SiteRepository) Get(ctx context.Context, id string) (*domain.Site, error) {
	query := `
		SELECT 
            s.site_id, s.tenant_id, s.site_name, s.location, 
            s.latitude, s.longitude, s.created_at, s.updated_at, t.tenant_name,
            (SELECT COUNT(*) FROM devices d WHERE d.site_id = s.site_id AND d.status != 'inactive') as device_count,
            (SELECT COUNT(*) FROM users u WHERE u.current_project_id IN (SELECT project_id FROM projects p WHERE p.site_id = s.site_id)) as worker_count
		FROM sites s
		LEFT JOIN tenants t ON s.tenant_id = t.tenant_id
		WHERE s.site_id = ?`

	var s domain.Site
	var tenantID, loc sql.NullString
	var lat, lng sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &tenantID, &s.Name, &loc,
		&lat, &lng, &s.CreatedAt, &s.UpdatedAt, &s.TenantName,
		&s.DeviceCount, &s.WorkerCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	if tenantID.Valid {
		s.TenantID = tenantID.String
	}
	if loc.Valid {
		s.Location = loc.String
	}
	if lat.Valid {
		s.Latitude = lat.Float64
	}
	if lng.Valid {
		s.Longitude = lng.Float64
	}

	return &s, nil
}

func (r *SiteRepository) List(ctx context.Context, tenantID string) ([]domain.Site, error) {
	query := `
        SELECT 
            s.site_id, s.tenant_id, s.site_name, s.location,
            s.latitude, s.longitude, s.created_at, s.updated_at, t.tenant_name,
            (SELECT COUNT(*) FROM devices d WHERE d.site_id = s.site_id AND d.status != 'inactive') as device_count,
            (SELECT COUNT(*) FROM users u WHERE u.current_project_id IN (SELECT project_id FROM projects p WHERE p.site_id = s.site_id)) as worker_count
        FROM sites s
        LEFT JOIN tenants t ON s.tenant_id = t.tenant_id
        WHERE 1=1 `

	args := []interface{}{}
	log.Printf("[SECURITY] SiteRepository.List: tenantID='%s'", tenantID)

	if tenantID != "" {
		query += " AND s.tenant_id = ?"
		args = append(args, tenantID)
	} else {
		log.Printf("[SECURITY] WARNING: Listing sites without tenantID filter")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []domain.Site
	for rows.Next() {
		var s domain.Site
		var tenantID, loc sql.NullString
		var lat, lng sql.NullFloat64
		if err := rows.Scan(
			&s.ID, &tenantID, &s.Name, &loc,
			&lat, &lng, &s.CreatedAt, &s.UpdatedAt, &s.TenantName,
			&s.DeviceCount, &s.WorkerCount,
		); err != nil {
			return nil, err
		}
		if tenantID.Valid {
			s.TenantID = tenantID.String
		}
		if loc.Valid {
			s.Location = loc.String
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
	query := `INSERT INTO sites (
        site_id, tenant_id, site_name, location, latitude, longitude
    ) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.TenantID, s.Name, s.Location, s.Latitude, s.Longitude,
	)
	return err
}

func (r *SiteRepository) Update(ctx context.Context, s *domain.Site) error {
	query := `UPDATE sites SET 
        tenant_id=?, site_name=?, location=?, latitude=?, longitude=? 
        WHERE site_id=?`
	_, err := r.db.ExecContext(ctx, query,
		s.TenantID, s.Name, s.Location, s.Latitude, s.Longitude, s.ID,
	)
	return err
}

func (r *SiteRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sites WHERE site_id = ?", id)
	return err
}
