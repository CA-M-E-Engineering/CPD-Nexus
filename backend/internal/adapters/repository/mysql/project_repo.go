package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/pkg/apperrors"
	"cpd-nexus/internal/pkg/idgen"
	"cpd-nexus/internal/pkg/logger"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ports.ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Get(ctx context.Context, userID, id string) (*domain.Project, error) {
	query := `
		SELECT 
            p.project_id, p.site_id, p.user_id, p.project_title, p.status, p.submission_entity,
            p.project_reference_number, p.project_contract_number, p.project_location_description, p.project_contract_name, p.hdb_precinct_name, 
            p.main_contractor_name, p.main_contractor_uen,
            p.worker_company_name, p.worker_company_uen,
            p.worker_company_client_name, p.worker_company_client_uen, p.worker_company_trade,
            p.offsite_fabricator_name, p.offsite_fabricator_uen, p.offsite_fabricator_location,
            p.pitstop_auth_id, pa.on_behalf_of_name as pitstop_auth_name,
            p.created_at, p.updated_at, s.site_name,
            (SELECT COUNT(*) FROM workers w WHERE w.current_project_id = p.project_id) as worker_count,
            (SELECT COUNT(*) FROM devices d WHERE d.site_id = p.site_id) as device_count
		FROM projects p
		LEFT JOIN sites s ON p.site_id = s.site_id
		LEFT JOIN pitstop_authorisations pa ON p.pitstop_auth_id = pa.pitstop_auth_id`

	args := []interface{}{id}
	if ports.IsVendor(ctx) {
		query += " WHERE p.project_id = ?"
	} else {
		query += " WHERE p.project_id = ? AND p.user_id = ?"
		args = append(args, userID)
	}

	var p domain.Project
	var siteID, scanUserID, status, ref, cRef, loc, cName, hdb sql.NullString
	var mcName, mcUEN, wcName, wcUEN, wccName, wccUEN, wcTrade, pitstopAuthID, pitstopAuthName sql.NullString
	var ofName, ofUEN, ofLoc sql.NullString

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&p.ID, &siteID, &scanUserID, &p.Title, &status, &p.SubmissionEntity,
		&ref, &cRef, &loc, &cName, &hdb,
		&mcName, &mcUEN, &wcName, &wcUEN, &wccName, &wccUEN, &wcTrade,
		&ofName, &ofUEN, &ofLoc,
		&pitstopAuthID, &pitstopAuthName,
		&p.CreatedAt, &p.UpdatedAt, &p.SiteName, &p.WorkerCount, &p.DeviceCount,
	)
	if err == sql.ErrNoRows {
		return nil, apperrors.NewNotFound("project", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	if siteID.Valid {
		p.SiteID = siteID.String
	}
	if scanUserID.Valid {
		p.UserID = scanUserID.String
	}
	if status.Valid {
		p.Status = status.String
	}
	if ref.Valid {
		p.Reference = ref.String
	}
	if cRef.Valid {
		p.ContractRef = cRef.String
	}
	if loc.Valid {
		p.Location = loc.String
	}
	if cName.Valid {
		p.ContractName = cName.String
	}
	if hdb.Valid {
		p.HDBPrecinct = hdb.String
	}
	if mcName.Valid {
		p.MainContractorName = mcName.String
	}
	if mcUEN.Valid {
		p.MainContractorUEN = mcUEN.String
	}
	if wcName.Valid {
		p.WorkerCompanyName = wcName.String
	}
	if wcUEN.Valid {
		p.WorkerCompanyUEN = wcUEN.String
	}
	if wccName.Valid {
		p.WorkerCompanyClientName = wccName.String
	}
	if wccUEN.Valid {
		p.WorkerCompanyClientUEN = wccUEN.String
	}
	if wcTrade.Valid {
		p.WorkerCompanyTrade = wcTrade.String
	}
	if pitstopAuthID.Valid {
		p.PitstopAuthID = &pitstopAuthID.String
	}
	if pitstopAuthName.Valid {
		p.PitstopAuthName = &pitstopAuthName.String
	}
	if ofName.Valid {
		p.OffsiteFabricatorName = ofName.String
	}
	if ofUEN.Valid {
		p.OffsiteFabricatorUEN = ofUEN.String
	}
	if ofLoc.Valid {
		p.OffsiteFabricatorLocation = ofLoc.String
	}

	return &p, nil
}

func (r *ProjectRepository) List(ctx context.Context, userID string) ([]domain.Project, error) {
	query := `
        SELECT 
            p.project_id, p.site_id, p.user_id, p.project_title, p.status, p.submission_entity,
            p.project_reference_number, p.project_contract_number, p.project_location_description, p.project_contract_name, p.hdb_precinct_name, 
            p.main_contractor_name, p.main_contractor_uen,
            p.worker_company_name, p.worker_company_uen,
            p.worker_company_client_name, p.worker_company_client_uen, p.worker_company_trade,
            p.offsite_fabricator_name, p.offsite_fabricator_uen, p.offsite_fabricator_location,
            p.pitstop_auth_id, pa.on_behalf_of_name as pitstop_auth_name,
            p.created_at, p.updated_at, s.site_name,
            (SELECT COUNT(*) FROM workers w WHERE w.current_project_id = p.project_id) as worker_count,
            (SELECT COUNT(*) FROM devices d WHERE d.site_id = p.site_id AND d.status != ?) as device_count
        FROM projects p
        LEFT JOIN sites s ON p.site_id = s.site_id
		LEFT JOIN pitstop_authorisations pa ON p.pitstop_auth_id = pa.pitstop_auth_id
        WHERE (p.status != ? OR p.status IS NULL)`

	if userID == "" && !ports.IsVendor(ctx) {
		return nil, apperrors.NewPermissionDenied("user_id is required for multi-tenant isolation")
	}
	args := []interface{}{}
	if !ports.IsVendor(ctx) {
		query += " AND p.user_id = ?"
		args = append(args, userID)
	}

	logger.Infof("[SECURITY] ProjectRepository.List: userID='%s'", userID)
	rows, err := r.db.QueryContext(ctx, query, append([]interface{}{domain.StatusInactive, domain.StatusInactive}, args...)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	rowCount := 0
	for rows.Next() {
		rowCount++
		var p domain.Project
		var siteID, uid, status, ref, cRef, loc, cName, hdb sql.NullString
		var mcName, mcUEN, wcName, wcUEN, wccName, wccUEN, wcTrade, pitstopAuthID, pitstopAuthName sql.NullString
		var ofName, ofUEN, ofLoc sql.NullString
		if err := rows.Scan(
			&p.ID, &siteID, &uid, &p.Title, &status, &p.SubmissionEntity,
			&ref, &cRef, &loc, &cName, &hdb,
			&mcName, &mcUEN, &wcName, &wcUEN, &wccName, &wccUEN, &wcTrade,
			&ofName, &ofUEN, &ofLoc,
			&pitstopAuthID, &pitstopAuthName,
			&p.CreatedAt, &p.UpdatedAt, &p.SiteName, &p.WorkerCount, &p.DeviceCount,
		); err != nil {
			logger.Infof("[ERROR] ProjectRepository.List Scan failed: %v", err)
			return nil, err
		}
		if siteID.Valid {
			p.SiteID = siteID.String
		}
		if uid.Valid {
			p.UserID = uid.String
		}
		if status.Valid {
			p.Status = status.String
		}
		if ref.Valid {
			p.Reference = ref.String
		}
		if cRef.Valid {
			p.ContractRef = cRef.String
		}
		if loc.Valid {
			p.Location = loc.String
		}
		if cName.Valid {
			p.ContractName = cName.String
		}
		if hdb.Valid {
			p.HDBPrecinct = hdb.String
		}
		if mcName.Valid {
			p.MainContractorName = mcName.String
		}
		if mcUEN.Valid {
			p.MainContractorUEN = mcUEN.String
		}
		if wcName.Valid {
			p.WorkerCompanyName = wcName.String
		}
		if wcUEN.Valid {
			p.WorkerCompanyUEN = wcUEN.String
		}
		if wccName.Valid {
			p.WorkerCompanyClientName = wccName.String
		}
		if wccUEN.Valid {
			p.WorkerCompanyClientUEN = wccUEN.String
		}
		if wcTrade.Valid {
			p.WorkerCompanyTrade = wcTrade.String
		}
		if pitstopAuthID.Valid {
			p.PitstopAuthID = &pitstopAuthID.String
		}
		if pitstopAuthName.Valid {
			p.PitstopAuthName = &pitstopAuthName.String
		}
		if ofName.Valid {
			p.OffsiteFabricatorName = ofName.String
		}
		if ofUEN.Valid {
			p.OffsiteFabricatorUEN = ofUEN.String
		}
		if ofLoc.Valid {
			p.OffsiteFabricatorLocation = ofLoc.String
		}

		projects = append(projects, p)
	}
	logger.Infof("[DEBUG] ProjectRepository.List: found %d projects for userID='%s'", len(projects), userID)
	return projects, nil
}

func (r *ProjectRepository) Create(ctx context.Context, p *domain.Project) error {
	id, err := idgen.GenerateNextID(r.db, "projects", "project_id", "project")
	if err != nil {
		return fmt.Errorf("failed to generate project ID: %w", err)
	}
	p.ID = id

	query := `INSERT INTO projects (
        project_id, site_id, user_id, project_title, status, submission_entity,
        project_reference_number, 
        project_contract_number, project_location_description, project_contract_name, hdb_precinct_name, 
        main_contractor_name, main_contractor_uen,
        worker_company_name, worker_company_uen,
        worker_company_client_name, worker_company_client_uen, worker_company_trade,
        offsite_fabricator_name, offsite_fabricator_uen, offsite_fabricator_location,
        pitstop_auth_id,
        created_at, updated_at
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err = r.db.ExecContext(ctx, query,
		p.ID, p.SiteID, p.UserID, p.Title, p.Status, p.SubmissionEntity, p.Reference,
		p.ContractRef, p.Location, p.ContractName, toNullString(p.HDBPrecinct),
		toNullString(p.MainContractorName), toNullString(p.MainContractorUEN),
		toNullString(p.WorkerCompanyName), toNullString(p.WorkerCompanyUEN),
		toNullString(p.WorkerCompanyClientName), toNullString(p.WorkerCompanyClientUEN), toNullString(p.WorkerCompanyTrade),
		toNullString(p.OffsiteFabricatorName), toNullString(p.OffsiteFabricatorUEN), toNullString(p.OffsiteFabricatorLocation),
		toNullableStringPtr(p.PitstopAuthID),
	)
	return err
}

func (r *ProjectRepository) Update(ctx context.Context, p *domain.Project) error {
	query := `UPDATE projects SET 
        site_id=?, user_id=?, project_title=?, status=?, submission_entity=?,
        project_reference_number=?, 
        project_contract_number=?, project_location_description=?, project_contract_name=?, hdb_precinct_name=?, 
        main_contractor_name=?, main_contractor_uen=?,
        worker_company_name=?, worker_company_uen=?,
        worker_company_client_name=?, worker_company_client_uen=?, worker_company_trade=?,
        offsite_fabricator_name=?, offsite_fabricator_uen=?, offsite_fabricator_location=?,
        pitstop_auth_id=?,
        updated_at=NOW()
        WHERE project_id=?`

	_, err := r.db.ExecContext(ctx, query,
		p.SiteID, p.UserID, p.Title, p.Status, p.SubmissionEntity, p.Reference,
		p.ContractRef, p.Location, p.ContractName, toNullString(p.HDBPrecinct),
		toNullString(p.MainContractorName), toNullString(p.MainContractorUEN),
		toNullString(p.WorkerCompanyName), toNullString(p.WorkerCompanyUEN),
		toNullString(p.WorkerCompanyClientName), toNullString(p.WorkerCompanyClientUEN), toNullString(p.WorkerCompanyTrade),
		toNullString(p.OffsiteFabricatorName), toNullString(p.OffsiteFabricatorUEN), toNullString(p.OffsiteFabricatorLocation),
		toNullableStringPtr(p.PitstopAuthID),
		p.ID,
	)
	return err
}

func toNullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func toNullableStringPtr(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}

func (r *ProjectRepository) Delete(ctx context.Context, userID, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Deactivate project
	res, err := tx.ExecContext(ctx, "UPDATE projects SET status = ? WHERE project_id = ? AND user_id = ?", domain.StatusInactive, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return apperrors.NewNotFound("project", id)
	}

	// 2. Unassign workers
	if _, err := tx.ExecContext(ctx, "UPDATE workers SET current_project_id = NULL WHERE current_project_id = ?", id); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ProjectRepository) AssignToSite(ctx context.Context, siteID string, projectIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Unassign all projects currently on this site
	_, err = tx.ExecContext(ctx, "UPDATE projects SET site_id = NULL WHERE site_id = ?", siteID)
	if err != nil {
		return fmt.Errorf("failed to clear old project assignments: %w", err)
	}

	// 2. Assign new projects
	if len(projectIDs) > 0 {
		stmt, err := tx.PrepareContext(ctx, "UPDATE projects SET site_id = ? WHERE project_id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, projId := range projectIDs {
			if _, err := stmt.ExecContext(ctx, siteID, projId); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
