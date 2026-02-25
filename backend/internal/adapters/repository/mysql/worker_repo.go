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
            w.worker_id, w.name, w.email, w.role, w.status, w.current_project_id,
            w.person_id_no, w.person_id_and_work_pass_type, w.person_nationality, w.person_trade, 
            p.project_title,
            s.site_name,
            s.location,
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
	var status, projID sql.NullString
	var pPassType, pNationality, pTrade sql.NullString
	var pName, sName, sLoc, uName, uID, uLat, uLng, uAdd sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID, &w.Name, &w.Email, &w.Role, &status, &projID,
		&w.PersonIDNo, &pPassType, &pNationality, &pTrade,
		&pName, &sName, &sLoc, &uName, &uID, &uLat, &uLng, &uAdd,
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
	if pPassType.Valid {
		w.PersonIDAndWorkPassType = pPassType.String
	}
	if pNationality.Valid {
		w.PersonNationality = pNationality.String
	}
	if pTrade.Valid {
		w.PersonTrade = pTrade.String
	}
	if pName.Valid {
		w.ProjectName = pName.String
	}
	if sName.Valid {
		w.SiteName = sName.String
	}
	if sLoc.Valid {
		w.SiteLocation = sLoc.String
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
            w.worker_id, w.name, w.email, w.role, w.status, w.current_project_id,
            w.person_id_no, w.person_id_and_work_pass_type, w.person_nationality, w.person_trade, 
            p.project_title,
            s.site_name,
            s.location,
            u.user_name,
            w.user_id,
            u.latitude,
            u.longitude,
            u.address
        FROM workers w
        LEFT JOIN projects p ON w.current_project_id = p.project_id
        LEFT JOIN sites s ON p.site_id = s.site_id
        LEFT JOIN users u ON w.user_id = u.user_id
        WHERE w.person_id_no = ? LIMIT 1`

	var w domain.Worker
	var status, projID sql.NullString
	var pPassType, pNationality, pTrade sql.NullString
	var pName, sName, sLoc, uName, uID, uLat, uLng, uAdd sql.NullString

	err := r.db.QueryRowContext(ctx, query, fin).Scan(
		&w.ID, &w.Name, &w.Email, &w.Role, &status, &projID,
		&w.PersonIDNo, &pPassType, &pNationality, &pTrade,
		&pName, &sName, &sLoc, &uName, &uID, &uLat, &uLng, &uAdd,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get worker by fin: %w", err)
	}
	if projID.Valid {
		w.CurrentProjectID = projID.String
	}
	if w.PersonIDNo != "" {
		// already scanned into struct
	}
	if pPassType.Valid {
		w.PersonIDAndWorkPassType = pPassType.String
	}
	if pNationality.Valid {
		w.PersonNationality = pNationality.String
	}
	if pTrade.Valid {
		w.PersonTrade = pTrade.String
	}
	if pTrade.Valid {
		w.PersonTrade = pTrade.String
	}
	if pName.Valid {
		w.ProjectName = pName.String
	}
	if sName.Valid {
		w.SiteName = sName.String
	}
	if sLoc.Valid {
		w.SiteLocation = sLoc.String
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
            w.worker_id, w.name, w.email, w.role, w.status, w.current_project_id,
            w.person_id_no, w.person_id_and_work_pass_type, w.person_nationality, w.person_trade, 
            p.project_title,
            s.site_name,
            s.location,
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
		var status, projID sql.NullString
		var pPassType, pNationality, pTrade sql.NullString
		var pName, sName, sLoc, uName, uID, uLat, uLng, uAdd sql.NullString

		if err := rows.Scan(
			&w.ID, &w.Name, &w.Email, &w.Role, &status, &projID,
			&w.PersonIDNo, &pPassType, &pNationality, &pTrade,
			&pName, &sName, &sLoc, &uName, &uID, &uLat, &uLng, &uAdd,
		); err != nil {
			log.Printf("[WorkerRepo] Scan error: %v", err)
			continue
		}

		if status.Valid {
			w.Status = status.String
		}
		if projID.Valid {
			w.CurrentProjectID = projID.String
		}
		if w.PersonIDNo != "" {
			// already scanned
		}
		if pPassType.Valid {
			w.PersonIDAndWorkPassType = pPassType.String
		}
		if pNationality.Valid {
			w.PersonNationality = pNationality.String
		}
		if pTrade.Valid {
			w.PersonTrade = pTrade.String
		}
		if pTrade.Valid {
			w.PersonTrade = pTrade.String
		}
		if pName.Valid {
			w.ProjectName = pName.String
		}
		if sName.Valid {
			w.SiteName = sName.String
		}
		if sLoc.Valid {
			w.SiteLocation = sLoc.String
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
        INSERT INTO workers (
            worker_id, user_id, name, email, role, status, current_project_id,
            person_id_no, person_id_and_work_pass_type, person_nationality, person_trade
        ) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	log.Printf("[WorkerRepo] Create: worker_id=%s user_id=%s name=%s role=%s", w.ID, w.UserID, w.Name, w.Role)

	_, err = r.db.ExecContext(ctx, query,
		w.ID, w.UserID, w.Name, w.Email, w.Role, w.Status,
		sql.NullString{String: w.CurrentProjectID, Valid: w.CurrentProjectID != ""},
		w.PersonIDNo, w.PersonIDAndWorkPassType, w.PersonNationality, w.PersonTrade,
	)
	if err != nil {
		return fmt.Errorf("failed to create worker: %w", err)
	}
	return nil
}

func (r *WorkerRepository) Update(ctx context.Context, w *domain.Worker) error {
	query := `
        UPDATE workers SET 
            name=?, email=?, status=?, role=?, current_project_id=?, user_id=?,
            person_id_no=?, person_id_and_work_pass_type=?, person_nationality=?, person_trade=?
        WHERE worker_id=?`

	_, err := r.db.ExecContext(ctx, query,
		w.Name, w.Email, w.Status, w.Role,
		sql.NullString{String: w.CurrentProjectID, Valid: w.CurrentProjectID != ""},
		w.UserID,
		w.PersonIDNo, w.PersonIDAndWorkPassType, w.PersonNationality, w.PersonTrade,
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
