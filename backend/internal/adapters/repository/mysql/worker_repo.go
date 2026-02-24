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

type WorkerRepository struct {
	db *sql.DB
}

func NewWorkerRepository(db *sql.DB) ports.WorkerRepository {
	return &WorkerRepository{db: db}
}

func (r *WorkerRepository) Get(ctx context.Context, id string) (*domain.Worker, error) {
	query := `
        SELECT 
            w.worker_id, w.name, w.email, w.role, w.status, w.trade_code, w.current_project_id, w.fin_nric, w.company_name,
            p.project_title,
            s.site_name,
            u.user_name,
            w.user_id,
            u.latitude,
            u.longitude,
            u.address
        FROM workers w
        LEFT JOIN projects p ON w.current_project_id = p.project_id
        LEFT JOIN sites s ON p.site_id = s.site_id
        LEFT JOIN users u ON w.user_id = u.user_id
        WHERE w.worker_id = ?`

	var w domain.Worker
	var status, tradeCode, projID, fin, cName sql.NullString
	var pName, sName, uName, uID, uLat, uLng, uAdd sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID, &w.Name, &w.Email, &w.Role, &status, &tradeCode, &projID, &fin, &cName,
		&pName, &sName, &uName, &uID, &uLat, &uLng, &uAdd,
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
	if uName.Valid {
		w.UserName = uName.String
	}
	if uID.Valid {
		w.UserID = uID.String
	}
	if uLat.Valid && uLng.Valid {
		w.UserLocation = uLat.String + ", " + uLng.String
	}
	if uAdd.Valid {
		w.UserAddress = uAdd.String
	}

	return &w, nil
}

func (r *WorkerRepository) GetByFIN(ctx context.Context, fin string) (*domain.Worker, error) {
	query := `
        SELECT 
            w.worker_id, w.name, w.email, w.role, w.status, w.trade_code, w.current_project_id, w.fin_nric, w.company_name,
            p.project_title,
            s.site_name,
            u.user_name,
            w.user_id,
            u.latitude,
            u.longitude,
            u.address
        FROM workers w
        LEFT JOIN projects p ON w.current_project_id = p.project_id
        LEFT JOIN sites s ON p.site_id = s.site_id
        LEFT JOIN users u ON w.user_id = u.user_id
        WHERE w.fin_nric = ? LIMIT 1`

	var w domain.Worker
	var status, tradeCode, projID, f, cName sql.NullString
	var pName, sName, uName, uID, uLat, uLng, uAdd sql.NullString

	err := r.db.QueryRowContext(ctx, query, fin).Scan(
		&w.ID, &w.Name, &w.Email, &w.Role, &status, &tradeCode, &projID, &f, &cName,
		&pName, &sName, &uName, &uID, &uLat, &uLng, &uAdd,
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
	if uName.Valid {
		w.UserName = uName.String
	}
	if uID.Valid {
		w.UserID = uID.String
	}
	if uLat.Valid && uLng.Valid {
		w.UserLocation = uLat.String + ", " + uLng.String
	}
	if uAdd.Valid {
		w.UserAddress = uAdd.String
	}

	return &w, nil
}

func (r *WorkerRepository) List(ctx context.Context, userID, siteID string) ([]domain.Worker, error) {
	query := `
        SELECT 
            w.worker_id, w.name, w.email, w.role, w.status, w.trade_code, w.current_project_id, w.fin_nric, w.company_name,
            p.project_title,
            s.site_name,
            u.user_name,
            w.user_id,
            u.latitude,
            u.longitude,
            u.address
        FROM workers w
        LEFT JOIN projects p ON w.current_project_id = p.project_id
        LEFT JOIN sites s ON p.site_id = s.site_id
        LEFT JOIN users u ON w.user_id = u.user_id
        WHERE w.role IN ('worker', 'pic', 'manager')`

	args := []interface{}{}
	if userID != "" {
		query += " AND w.user_id = ?"
		args = append(args, userID)
	}
	if siteID != "" {
		query += " AND s.site_id = ?"
		args = append(args, siteID)
	}

	query += " ORDER BY w.role DESC, w.name ASC"

	log.Printf("[WorkerRepo] List: userID=%s, siteID=%s", userID, siteID)
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
		var pName, sName, uName, uID, uLat, uLng, uAdd sql.NullString

		if err := rows.Scan(
			&w.ID, &w.Name, &w.Email, &w.Role, &status, &tradeCode, &projID, &fin, &cName,
			&pName, &sName, &uName, &uID, &uLat, &uLng, &uAdd,
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
		if uName.Valid {
			w.UserName = uName.String
		}
		if uID.Valid {
			w.UserID = uID.String
		}
		if uLat.Valid && uLng.Valid {
			w.UserLocation = uLat.String + ", " + uLng.String
		}
		if uAdd.Valid {
			w.UserAddress = uAdd.String
		}

		workers = append(workers, w)
	}
	return workers, nil
}

func (r *WorkerRepository) Create(ctx context.Context, w *domain.Worker) error {
	id, err := idgen.GenerateNextID(r.db, "workers", "worker_id", "worker")
	if err != nil {
		return fmt.Errorf("failed to generate worker ID: %w", err)
	}
	w.ID = id

	query := `
        INSERT INTO workers (worker_id, user_id, name, email, role, status, trade_code, current_project_id, fin_nric, company_name) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = r.db.ExecContext(ctx, query,
		w.ID, w.UserID, w.Name, w.Email, w.Role, w.Status, w.TradeCode,
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
	query := "UPDATE workers SET name=?, email=?, trade_code=?, status=?, role=?, current_project_id=?, fin_nric=?, user_id=?, company_name=? WHERE worker_id=?"

	_, err := r.db.ExecContext(ctx, query,
		w.Name, w.Email, w.TradeCode, w.Status, w.Role,
		sql.NullString{String: w.CurrentProjectID, Valid: w.CurrentProjectID != ""},
		w.FIN, w.UserID,
		sql.NullString{String: w.CompanyName, Valid: w.CompanyName != ""},
		w.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update worker: %w", err)
	}
	return nil
}

func (r *WorkerRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE workers SET status = 'inactive' WHERE worker_id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to deactivate worker: %w", err)
	}
	return nil
}

func (r *WorkerRepository) GetProjectUserID(ctx context.Context, projectID string) (string, error) {
	var projectUserID string
	err := r.db.QueryRowContext(ctx, "SELECT user_id FROM projects WHERE project_id = ?", projectID).Scan(&projectUserID)
	if err != nil {
		return "", err
	}
	return projectUserID, nil
}
