package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) ports.TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) GetByUsername(ctx context.Context, username string) (*domain.Tenant, error) {
	query := `SELECT tenant_id, tenant_name, username, tenant_type, status, latitude, longitude FROM tenants WHERE username = ? AND status = 'active'`

	var t domain.Tenant
	var lat, lng sql.NullString
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&t.ID, &t.Name, &t.Username, &t.TenantType, &t.Status, &lat, &lng,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by username: %w", err)
	}

	if lat.Valid {
		t.Latitude = lat.String
	}
	if lng.Valid {
		t.Longitude = lng.String
	}

	return &t, nil
}

func (r *TenantRepository) Get(ctx context.Context, id string) (*domain.Tenant, error) {
	query := `SELECT tenant_id, tenant_name, username, tenant_type, status, latitude, longitude FROM tenants WHERE tenant_id = ?`

	var t domain.Tenant
	var lat, lng sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.Username, &t.TenantType, &t.Status, &lat, &lng,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	if lat.Valid {
		t.Latitude = lat.String
	}
	if lng.Valid {
		t.Longitude = lng.String
	}

	return &t, nil
}
