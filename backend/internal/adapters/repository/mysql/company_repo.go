package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type CompanyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) ports.CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Get(ctx context.Context, id string) (*domain.Company, error) {
	query := `SELECT company_id, tenant_id, company_name, uen, company_type, address, latitude, longitude, status, created_at, updated_at 
			  FROM companies WHERE company_id = ?`

	var c domain.Company
	var addr sql.NullString
	var lat, lng sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.TenantID, &c.Name, &c.UEN, &c.CompanyType, &addr, &lat, &lng, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	c.Address = addr.String
	c.Latitude = lat.Float64
	c.Longitude = lng.Float64

	return &c, nil
}

func (r *CompanyRepository) GetByUEN(ctx context.Context, uen string) (*domain.Company, error) {
	query := `SELECT company_id, tenant_id, company_name, uen, company_type, address, latitude, longitude, status, created_at, updated_at 
			  FROM companies WHERE uen = ?`

	var c domain.Company
	var addr sql.NullString
	var lat, lng sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, uen).Scan(
		&c.ID, &c.TenantID, &c.Name, &c.UEN, &c.CompanyType, &addr, &lat, &lng, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get company by UEN: %w", err)
	}

	c.Address = addr.String
	c.Latitude = lat.Float64
	c.Longitude = lng.Float64

	return &c, nil
}

func (r *CompanyRepository) ListByTenant(ctx context.Context, tenantID string) ([]domain.Company, error) {
	query := `SELECT company_id, tenant_id, company_name, uen, company_type, address, latitude, longitude, status, created_at, updated_at 
			  FROM companies WHERE tenant_id = ?`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}
	defer rows.Close()

	var results []domain.Company
	for rows.Next() {
		var c domain.Company
		var addr sql.NullString
		var lat, lng sql.NullFloat64
		err := rows.Scan(
			&c.ID, &c.TenantID, &c.Name, &c.UEN, &c.CompanyType, &addr, &lat, &lng, &c.Status, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}
		c.Address = addr.String
		c.Latitude = lat.Float64
		c.Longitude = lng.Float64
		results = append(results, c)
	}

	return results, nil
}

func (r *CompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	query := `INSERT INTO companies (company_id, tenant_id, company_name, uen, company_type, address, latitude, longitude, status)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	if company.ID == "" {
		return fmt.Errorf("failed to create company: company_id is required")
	}

	_, err := r.db.ExecContext(ctx, query,
		company.ID, company.TenantID, company.Name, company.UEN, company.CompanyType, company.Address, company.Latitude, company.Longitude, company.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

func (r *CompanyRepository) Update(ctx context.Context, company *domain.Company) error {
	query := `UPDATE companies SET tenant_id = ?, company_name = ?, uen = ?, company_type = ?, address = ?, latitude = ?, longitude = ?, status = ?
			  WHERE company_id = ?`

	_, err := r.db.ExecContext(ctx, query,
		company.TenantID, company.Name, company.UEN, company.CompanyType, company.Address, company.Latitude, company.Longitude, company.Status, company.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}
	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM companies WHERE company_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}
	return nil
}
