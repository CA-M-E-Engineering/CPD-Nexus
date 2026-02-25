package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"sgbuildex/internal/pkg/idgen"
	"strings"
)

type Scanner interface {
	Scan(dest ...interface{}) error
}

type WorkerRepository struct {
	db *sql.DB
}

func NewWorkerRepository(db *sql.DB) ports.WorkerRepository {
	return &WorkerRepository{db: db}
}

func (r *WorkerRepository) Get(ctx context.Context, id string) (*domain.Worker, error) {
	query := workerBaseSelect + " WHERE w.worker_id = ?"
	return r.scanRow(r.db.QueryRowContext(ctx, query, id))
}

func (r *WorkerRepository) GetByFIN(ctx context.Context, fin string) (*domain.Worker, error) {
	query := workerBaseSelect + " WHERE w.person_id_no = ? LIMIT 1"
	return r.scanRow(r.db.QueryRowContext(ctx, query, fin))
}

const workerBaseSelect = `
    SELECT 
        w.worker_id, w.name, w.email, w.role, w.user_type, w.status, w.current_project_id,
        w.person_id_no, w.person_id_and_work_pass_type, w.person_nationality, w.person_trade, 
        w.auth_start_time, w.auth_end_time, w.fdid, w.face_img_loc, w.card_number, w.card_type, w.is_synced,
        p.project_title,
        s.site_name,
        s.location,
        u.user_name,
        w.user_id,
        u.latitude,
        u.longitude,
        u.address,
        p.site_id
    FROM workers w
    LEFT JOIN projects p ON w.current_project_id = p.project_id
    LEFT JOIN sites s ON p.site_id = s.site_id
    LEFT JOIN users u ON w.user_id = u.user_id`

func (r *WorkerRepository) List(ctx context.Context, userID, siteID string) ([]domain.Worker, error) {
	query := workerBaseSelect + " WHERE w.status = '" + domain.StatusActive + "' AND w.role IN ('worker', 'pic', 'manager')"

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
		w, err := r.scanRow(rows)
		if err != nil {
			log.Printf("[WorkerRepo] Scan error: %v", err)
			continue
		}
		workers = append(workers, *w)
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
            worker_id, user_id, name, email, role, user_type, status, current_project_id,
            person_id_no, person_id_and_work_pass_type, person_nationality, person_trade,
            auth_start_time, auth_end_time, fdid, face_img_loc, card_number, card_type, is_synced
        ) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	log.Printf("[WorkerRepo] Create: worker_id=%s user_id=%s name=%s role=%s", w.ID, w.UserID, w.Name, w.Role)

	// DB constraint: fdid defaults to 1
	fdidToInsert := w.FDID
	if fdidToInsert == 0 {
		fdidToInsert = 1
	}

	userType := w.UserType
	if userType == "" {
		userType = domain.UserTypeUser
	}

	_, err = r.db.ExecContext(ctx, query,
		w.ID, w.UserID, w.Name, w.Email, w.Role, userType, w.Status,
		sql.NullString{String: w.CurrentProjectID, Valid: w.CurrentProjectID != ""},
		w.PersonIDNo, w.PersonIDAndWorkPassType, w.PersonNationality, w.PersonTrade,
		sql.NullString{String: formatMySQLDate(w.AuthStartTime), Valid: w.AuthStartTime != ""},
		sql.NullString{String: formatMySQLDate(w.AuthEndTime), Valid: w.AuthEndTime != ""},
		fdidToInsert,
		sql.NullString{String: w.FaceImgLoc, Valid: w.FaceImgLoc != ""},
		sql.NullString{String: w.CardNumber, Valid: w.CardNumber != ""},
		sql.NullString{String: w.CardType, Valid: w.CardType != ""},
		w.IsSynced,
	)
	if err != nil {
		return fmt.Errorf("failed to create worker: %w", err)
	}
	return nil
}

func (r *WorkerRepository) Update(ctx context.Context, w *domain.Worker) error {
	query := `
        UPDATE workers SET 
            name=?, email=?, status=?, role=?, user_type=?, current_project_id=?, user_id=?,
            person_id_no=?, person_id_and_work_pass_type=?, person_nationality=?, person_trade=?,
            auth_start_time=?, auth_end_time=?, fdid=?, face_img_loc=?, card_number=?, card_type=?, is_synced=?
        WHERE worker_id=?`

	fdidToInsert := w.FDID
	if fdidToInsert == 0 {
		fdidToInsert = 1
	}

	userType := w.UserType
	if userType == "" {
		userType = "user"
	}

	_, err := r.db.ExecContext(ctx, query,
		w.Name, w.Email, w.Status, w.Role, userType,
		sql.NullString{String: w.CurrentProjectID, Valid: w.CurrentProjectID != ""},
		w.UserID,
		w.PersonIDNo, w.PersonIDAndWorkPassType, w.PersonNationality, w.PersonTrade,
		sql.NullString{String: formatMySQLDate(w.AuthStartTime), Valid: w.AuthStartTime != ""},
		sql.NullString{String: formatMySQLDate(w.AuthEndTime), Valid: w.AuthEndTime != ""},
		fdidToInsert,
		sql.NullString{String: w.FaceImgLoc, Valid: w.FaceImgLoc != ""},
		sql.NullString{String: w.CardNumber, Valid: w.CardNumber != ""},
		sql.NullString{String: w.CardType, Valid: w.CardType != ""},
		w.IsSynced,
		w.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update worker: %w", err)
	}
	return nil
}

func (r *WorkerRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE workers SET status = ?, current_project_id = NULL WHERE worker_id = ?", domain.StatusInactive, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate worker and clear project: %w", err)
	}
	return nil
}

func (r *WorkerRepository) ListByIsSynced(ctx context.Context, userID string, syncStatus int) ([]domain.Worker, error) {
	query := workerBaseSelect + " WHERE w.is_synced = ? AND w.status = '" + domain.StatusActive + "'"

	args := []interface{}{syncStatus}
	if userID != "" {
		query += " AND w.user_id = ?"
		args = append(args, userID)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list workers by sync status: %w", err)
	}
	defer rows.Close()

	var workers []domain.Worker
	for rows.Next() {
		w, err := r.scanRow(rows)
		if err != nil {
			log.Printf("[WorkerRepo] ListByIsSynced scan error: %v", err)
			continue
		}
		workers = append(workers, *w)
	}
	return workers, nil
}

func (r *WorkerRepository) scanRow(scanner Scanner) (*domain.Worker, error) {
	var w domain.Worker
	var status, projID, siteID sql.NullString
	var userType sql.NullString
	var pPassType, pNationality, pTrade sql.NullString
	var pName, sName, sLoc, uName, uID, uLat, uLng, uAdd sql.NullString
	var aStart, aEnd, fImg, cNum, cType sql.NullString
	var fdid, isSynced sql.NullInt64

	err := scanner.Scan(
		&w.ID, &w.Name, &w.Email, &w.Role, &userType, &status, &projID,
		&w.PersonIDNo, &pPassType, &pNationality, &pTrade,
		&aStart, &aEnd, &fdid, &fImg, &cNum, &cType, &isSynced,
		&pName, &sName, &sLoc, &uName, &uID, &uLat, &uLng, &uAdd,
		&siteID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if userType.Valid {
		w.UserType = userType.String
	}
	if status.Valid {
		w.Status = status.String
	}
	if projID.Valid {
		w.CurrentProjectID = projID.String
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
	if aStart.Valid {
		w.AuthStartTime = formatMySQLDate(aStart.String)
	}
	if aEnd.Valid {
		w.AuthEndTime = formatMySQLDate(aEnd.String)
	}
	if fdid.Valid {
		w.FDID = int(fdid.Int64)
	}
	if isSynced.Valid {
		w.IsSynced = int(isSynced.Int64)
	}
	if fImg.Valid {
		w.FaceImgLoc = fImg.String
	}
	if cNum.Valid {
		w.CardNumber = cNum.String
	}
	if cType.Valid {
		w.CardType = cType.String
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
	if siteID.Valid {
		w.SiteID = siteID.String
	}
	if uLat.Valid && uLng.Valid {
		w.UserLocation = uLat.String + ", " + uLng.String
	}
	if uAdd.Valid {
		w.UserAddress = uAdd.String
	}

	return &w, nil
}

func (r *WorkerRepository) AssignToProject(ctx context.Context, projectID string, workerIDs []string, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Unassign all workers currently on this project for this user
	_, err = tx.ExecContext(ctx, "UPDATE workers SET current_project_id = NULL WHERE current_project_id = ? AND user_id = ?", projectID, userID)
	if err != nil {
		return fmt.Errorf("failed to clear old assignments: %w", err)
	}

	// 2. Assign new workers (only if they are active)
	if len(workerIDs) > 0 {
		query := "UPDATE workers SET current_project_id = ? WHERE worker_id = ? AND user_id = ? AND status = ?"
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, wid := range workerIDs {
			_, err = stmt.ExecContext(ctx, projectID, wid, userID, domain.StatusActive)
			if err != nil {
				return fmt.Errorf("failed to assign worker %s: %w", wid, err)
			}
		}
	}

	return tx.Commit()
}

func (r *WorkerRepository) GetProjectUserID(ctx context.Context, projectID string) (string, error) {
	var projectUserID string
	err := r.db.QueryRowContext(ctx, "SELECT user_id FROM projects WHERE project_id = ?", projectID).Scan(&projectUserID)
	if err != nil {
		return "", err
	}
	return projectUserID, nil
}

func formatMySQLDate(ds string) string {
	if len(ds) >= 19 {
		tmp := strings.Replace(ds, "T", " ", 1)
		tmp = strings.Replace(tmp, "Z", "", 1)
		return tmp
	}
	return ds
}
