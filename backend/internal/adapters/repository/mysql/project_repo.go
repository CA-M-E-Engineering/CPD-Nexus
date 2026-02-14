package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ports.ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Get(ctx context.Context, id string) (*domain.Project, error) {
	query := `
		SELECT 
            p.project_id, p.site_id, p.tenant_id, p.project_title, p.status, 
            p.project_reference_number, p.project_contract_number, p.project_location_description, p.project_contract_name, p.hdb_precinct_name, 
            p.main_contractor_id, p.offsite_fabricator_id, p.worker_company_id, p.worker_company_client_id,
            c1.company_name as mc_name, c2.company_name as of_name, c3.company_name as wc_name, c4.company_name as wcc_name,
            p.created_at, p.updated_at, s.site_name,
            (SELECT COUNT(*) FROM users u WHERE u.current_project_id = p.project_id) as worker_count,
            (SELECT COUNT(*) FROM devices d WHERE d.site_id = p.site_id) as device_count
		FROM projects p
		LEFT JOIN sites s ON p.site_id = s.site_id
        LEFT JOIN companies c1 ON p.main_contractor_id = c1.company_id
        LEFT JOIN companies c2 ON p.offsite_fabricator_id = c2.company_id
        LEFT JOIN companies c3 ON p.worker_company_id = c3.company_id
        LEFT JOIN companies c4 ON p.worker_company_client_id = c4.company_id
		WHERE p.project_id = ?`

	var p domain.Project
	var siteID, tenantID, status, ref, cRef, loc, cName, hdb, mcID, ofID, wcID, wccID sql.NullString
	var mcn, ofn, wcn, wccn sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &siteID, &tenantID, &p.Title, &status,
		&ref, &cRef, &loc, &cName, &hdb, &mcID, &ofID, &wcID, &wccID,
		&mcn, &ofn, &wcn, &wccn,
		&p.CreatedAt, &p.UpdatedAt, &p.SiteName, &p.WorkerCount, &p.DeviceCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	if siteID.Valid {
		p.SiteID = siteID.String
	}
	if tenantID.Valid {
		p.TenantID = tenantID.String
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
	if mcID.Valid {
		p.MainContractorID = mcID.String
	}
	if ofID.Valid {
		p.OffsiteFabricatorID = ofID.String
	}
	if wcID.Valid {
		p.WorkerCompanyID = wcID.String
	}
	if wccID.Valid {
		p.WorkerCompanyClientID = wccID.String
	}
	if mcn.Valid {
		p.MainContractorName = mcn.String
	}
	if ofn.Valid {
		p.OffsiteFabricatorName = ofn.String
	}
	if wcn.Valid {
		p.WorkerCompanyName = wcn.String
	}
	if wccn.Valid {
		p.WorkerCompanyClientName = wccn.String
	}

	return &p, nil
}

func (r *ProjectRepository) List(ctx context.Context, tenantID string) ([]domain.Project, error) {
	query := `
        SELECT 
            p.project_id, p.site_id, p.tenant_id, p.project_title, p.status, 
            p.project_reference_number, p.project_contract_number, p.project_location_description, p.project_contract_name, p.hdb_precinct_name, 
            p.main_contractor_id, p.offsite_fabricator_id, p.worker_company_id, p.worker_company_client_id,
            c1.company_name as mc_name, c2.company_name as of_name, c3.company_name as wc_name, c4.company_name as wcc_name,
            p.created_at, p.updated_at, s.site_name,
            (SELECT COUNT(*) FROM users u WHERE u.current_project_id = p.project_id) as worker_count,
            (SELECT COUNT(*) FROM devices d WHERE d.site_id = p.site_id AND d.status != 'inactive') as device_count
        FROM projects p
        LEFT JOIN sites s ON p.site_id = s.site_id
        LEFT JOIN companies c1 ON p.main_contractor_id = c1.company_id
        LEFT JOIN companies c2 ON p.offsite_fabricator_id = c2.company_id
        LEFT JOIN companies c3 ON p.worker_company_id = c3.company_id
        LEFT JOIN companies c4 ON p.worker_company_client_id = c4.company_id
        WHERE (p.status != 'inactive' OR p.status IS NULL)`

	args := []interface{}{}
	if tenantID != "" {
		query += " AND p.tenant_id = ?"
		args = append(args, tenantID)
	}

	log.Printf("[SECURITY] ProjectRepository.List: tenantID='%s'", tenantID)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		var siteID, tenantID, status, ref, cRef, loc, cName, hdb, mcID, ofID, wcID, wccID sql.NullString
		var mcn, ofn, wcn, wccn sql.NullString
		if err := rows.Scan(
			&p.ID, &siteID, &tenantID, &p.Title, &status,
			&ref, &cRef, &loc, &cName, &hdb, &mcID, &ofID, &wcID, &wccID,
			&mcn, &ofn, &wcn, &wccn,
			&p.CreatedAt, &p.UpdatedAt, &p.SiteName, &p.WorkerCount, &p.DeviceCount,
		); err != nil {
			return nil, err
		}
		if siteID.Valid {
			p.SiteID = siteID.String
		}
		if tenantID.Valid {
			p.TenantID = tenantID.String
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
		if mcID.Valid {
			p.MainContractorID = mcID.String
		}
		if ofID.Valid {
			p.OffsiteFabricatorID = ofID.String
		}
		if wcID.Valid {
			p.WorkerCompanyID = wcID.String
		}
		if wccID.Valid {
			p.WorkerCompanyClientID = wccID.String
		}
		if mcn.Valid {
			p.MainContractorName = mcn.String
		}
		if ofn.Valid {
			p.OffsiteFabricatorName = ofn.String
		}
		if wcn.Valid {
			p.WorkerCompanyName = wcn.String
		}
		if wccn.Valid {
			p.WorkerCompanyClientName = wccn.String
		}

		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepository) Create(ctx context.Context, p *domain.Project) error {
	query := `INSERT INTO projects (
        project_id, site_id, tenant_id, project_title, status, project_reference_number, 
        project_contract_number, project_location_description, project_contract_name, hdb_precinct_name, 
        main_contractor_id, offsite_fabricator_id, worker_company_id, worker_company_client_id,
        created_at, updated_at
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.SiteID, p.TenantID, p.Title, p.Status, p.Reference,
		p.ContractRef, p.Location, p.ContractName, toNullString(p.HDBPrecinct),
		toNullString(p.MainContractorID), toNullString(p.OffsiteFabricatorID),
		toNullString(p.WorkerCompanyID), toNullString(p.WorkerCompanyClientID),
	)
	return err
}

func (r *ProjectRepository) Update(ctx context.Context, p *domain.Project) error {
	query := `UPDATE projects SET 
        site_id=?, tenant_id=?, project_title=?, status=?, project_reference_number=?, 
        project_contract_number=?, project_location_description=?, project_contract_name=?, hdb_precinct_name=?, 
        main_contractor_id=?, offsite_fabricator_id=?, worker_company_id=?, worker_company_client_id=?,
        updated_at=NOW()
        WHERE project_id=?`

	_, err := r.db.ExecContext(ctx, query,
		p.SiteID, p.TenantID, p.Title, p.Status, p.Reference,
		p.ContractRef, p.Location, p.ContractName, toNullString(p.HDBPrecinct),
		toNullString(p.MainContractorID), toNullString(p.OffsiteFabricatorID),
		toNullString(p.WorkerCompanyID), toNullString(p.WorkerCompanyClientID),
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

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE projects SET status = 'inactive' WHERE project_id = ?", id)
	return err
}
