package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"

	"github.com/google/uuid"
)

type WorkerRepository struct {
	db *sql.DB
}

func NewWorkerRepository(db *sql.DB) ports.WorkerRepository {
	return &WorkerRepository{db: db}
}

func (r *WorkerRepository) Get(ctx context.Context, id string) (*domain.Worker, error) {
	query := `
        SELECT 
            u.user_id, u.name, u.email, u.role, u.status, u.trade_code, u.current_project_id, u.fin_nric, u.company_name,
            p.project_title,
            s.site_name,
            t.tenant_name,
            u.tenant_id,
            t.latitude,
            t.longitude,
            t.address
        FROM users u
        LEFT JOIN projects p ON u.current_project_id = p.project_id
        LEFT JOIN sites s ON p.site_id = s.site_id
        LEFT JOIN tenants t ON u.tenant_id = t.tenant_id
        WHERE u.user_id = ?`

	var w domain.Worker
	var status, tradeCode, projID, fin, cName sql.NullString
	var pName, sName, tName, tID, tLat, tLng, tAdd sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID, &w.Name, &w.Email, &w.Role, &status, &tradeCode, &projID, &fin, &cName,
		&pName, &sName, &tName, &tID, &tLat, &tLng, &tAdd,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}

	if status.Valid {
		w.Status = status.String
	}
	if tradeCode.Valid {
		w.TradeCode = tradeCode.String
	}
	if projID.Valid {
		w.CurrentProjectID = projID.String
	}
	if fin.Valid {
		w.FIN = fin.String
	}
	if cName.Valid {
		w.CompanyName = cName.String
	}
	if pName.Valid {
		w.ProjectName = pName.String
	}
	if sName.Valid {
		w.SiteName = sName.String
	}
	if tName.Valid {
		w.TenantName = tName.String
	}
	if tID.Valid {
		w.TenantID = tID.String
	}
	if tLat.Valid && tLng.Valid {
		w.TenantLocation = tLat.String + ", " + tLng.String
	}
	if tAdd.Valid {
		w.TenantAddress = tAdd.String
	}

	return &w, nil
}

func (r *WorkerRepository) GetByFIN(ctx context.Context, fin string) (*domain.Worker, error) {
	query := `
        SELECT 
            u.user_id, u.name, u.email, u.role, u.status, u.trade_code, u.current_project_id, u.fin_nric, u.company_name,
            p.project_title,
            s.site_name,
            t.tenant_name,
            u.tenant_id,
            t.latitude,
            t.longitude,
            t.address
        FROM users u
        LEFT JOIN projects p ON u.current_project_id = p.project_id
        LEFT JOIN sites s ON p.site_id = s.site_id
        LEFT JOIN tenants t ON u.tenant_id = t.tenant_id
        WHERE u.fin_nric = ? LIMIT 1`

	var w domain.Worker
	var status, tradeCode, projID, f, cName sql.NullString
	var pName, sName, tName, tID, tLat, tLng, tAdd sql.NullString

	err := r.db.QueryRowContext(ctx, query, fin).Scan(
		&w.ID, &w.Name, &w.Email, &w.Role, &status, &tradeCode, &projID, &f, &cName,
		&pName, &sName, &tName, &tID, &tLat, &tLng, &tAdd,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get worker by fin: %w", err)
	}

	if status.Valid {
		w.Status = status.String
	}
	if tradeCode.Valid {
		w.TradeCode = tradeCode.String
	}
	if projID.Valid {
		w.CurrentProjectID = projID.String
	}
	if f.Valid {
		w.FIN = f.String
	}
	if cName.Valid {
		w.CompanyName = cName.String
	}
	if pName.Valid {
		w.ProjectName = pName.String
	}
	if sName.Valid {
		w.SiteName = sName.String
	}
	if tName.Valid {
		w.TenantName = tName.String
	}
	if tID.Valid {
		w.TenantID = tID.String
	}
	if tLat.Valid && tLng.Valid {
		w.TenantLocation = tLat.String + ", " + tLng.String
	}
	if tAdd.Valid {
		w.TenantAddress = tAdd.String
	}

	return &w, nil
}

func (r *WorkerRepository) List(ctx context.Context, tenantID, siteID string) ([]domain.Worker, error) {
	query := `
        SELECT 
            u.user_id, u.name, u.email, u.role, u.status, u.trade_code, u.current_project_id, u.fin_nric, u.company_name,
            p.project_title,
            s.site_name,
            t.tenant_name,
            u.tenant_id,
            t.latitude,
            t.longitude,
            t.address
        FROM users u
        LEFT JOIN projects p ON u.current_project_id = p.project_id
        LEFT JOIN sites s ON p.site_id = s.site_id
        LEFT JOIN tenants t ON u.tenant_id = t.tenant_id
        WHERE u.role IN ('worker', 'pic', 'manager') AND (u.status != 'inactive' OR u.status IS NULL)`

	args := []interface{}{}
	if tenantID != "" {
		query += " AND u.tenant_id = ?"
		args = append(args, tenantID)
	}
	if siteID != "" {
		query += " AND s.site_id = ?"
		args = append(args, siteID)
	}

	query += " ORDER BY u.role DESC, u.name ASC"

	log.Printf("[WorkerRepo] List: tenantID=%s, siteID=%s", tenantID, siteID)
	log.Printf("[WorkerRepo] Executing query: %s", query)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list workers: %w", err)
	}
	defer rows.Close()

	var workers []domain.Worker
	for rows.Next() {
		var w domain.Worker
		var status, tradeCode, projID, fin, cName sql.NullString
		var pName, sName, tName, tID, tLat, tLng, tAdd sql.NullString

		if err := rows.Scan(
			&w.ID, &w.Name, &w.Email, &w.Role, &status, &tradeCode, &projID, &fin, &cName,
			&pName, &sName, &tName, &tID, &tLat, &tLng, &tAdd,
		); err != nil {
			log.Printf("[WorkerRepo] Scan error: %v", err)
			continue
		}

		if status.Valid {
			w.Status = status.String
		}
		if tradeCode.Valid {
			w.TradeCode = tradeCode.String
		}
		if projID.Valid {
			w.CurrentProjectID = projID.String
		}
		if fin.Valid {
			w.FIN = fin.String
		}
		if cName.Valid {
			w.CompanyName = cName.String
		}
		if pName.Valid {
			w.ProjectName = pName.String
		}
		if sName.Valid {
			w.SiteName = sName.String
		}
		if tName.Valid {
			w.TenantName = tName.String
		}
		if tID.Valid {
			w.TenantID = tID.String
		}
		if tLat.Valid && tLng.Valid {
			w.TenantLocation = tLat.String + ", " + tLng.String
		}
		if tAdd.Valid {
			w.TenantAddress = tAdd.String
		}

		workers = append(workers, w)
	}
	return workers, nil
}

func (r *WorkerRepository) Create(ctx context.Context, w *domain.Worker) error {
	if w.ID == "" {
		w.ID = "worker-" + uuid.New().String()
	}

	query := `
        INSERT INTO users (user_id, tenant_id, name, email, role, status, trade_code, current_project_id, fin_nric, company_name) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		w.ID, w.TenantID, w.Name, w.Email, w.Role, w.Status, w.TradeCode,
		sql.NullString{String: w.CurrentProjectID, Valid: w.CurrentProjectID != ""},
		w.FIN,
		sql.NullString{String: w.CompanyName, Valid: w.CompanyName != ""},
	)
	if err != nil {
		return fmt.Errorf("failed to create worker: %w", err)
	}
	return nil
}

func (r *WorkerRepository) Update(ctx context.Context, w *domain.Worker) error {
	query := "UPDATE users SET name=?, email=?, trade_code=?, status=?, role=?, current_project_id=?, fin_nric=?, tenant_id=?, company_name=? WHERE user_id=?"

	_, err := r.db.ExecContext(ctx, query,
		w.Name, w.Email, w.TradeCode, w.Status, w.Role,
		sql.NullString{String: w.CurrentProjectID, Valid: w.CurrentProjectID != ""},
		w.FIN, w.TenantID,
		sql.NullString{String: w.CompanyName, Valid: w.CompanyName != ""},
		w.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update worker: %w", err)
	}
	return nil
}

func (r *WorkerRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET status = 'inactive' WHERE user_id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to deactivate worker: %w", err)
	}
	return nil
}

func (r *WorkerRepository) GetProjectTenantID(ctx context.Context, projectID string) (string, error) {
	var projectTenantID string
	err := r.db.QueryRowContext(ctx, "SELECT tenant_id FROM projects WHERE project_id = ?", projectID).Scan(&projectTenantID)
	if err != nil {
		return "", err
	}
	return projectTenantID, nil
}
