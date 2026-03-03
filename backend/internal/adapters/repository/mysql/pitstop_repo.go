package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sgbuildex/internal/core/domain"
	"strings"
)

type PitstopRepository struct {
	db *sql.DB
}

func NewPitstopRepository(db *sql.DB) *PitstopRepository {
	return &PitstopRepository{db: db}
}

func (r *PitstopRepository) GetAuthorisations(ctx context.Context) ([]*domain.PitstopAuthorisation, error) {
	query := `
		SELECT pitstop_auth_id, dataset_id, dataset_name, user_id, 
		       regulator_id, regulator_name, maincon_id, maincon_name, status, last_synced_at
		FROM pitstop_authorisations
		ORDER BY dataset_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pitstop authorisations: %w", err)
	}
	defer rows.Close()

	var auths []*domain.PitstopAuthorisation
	for rows.Next() {
		var a domain.PitstopAuthorisation
		if err := rows.Scan(
			&a.PitstopAuthID, &a.DatasetID, &a.DatasetName, &a.UserID,
			&a.RegulatorID, &a.RegulatorName, &a.MainconID, &a.MainconName, &a.Status, &a.LastSyncedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan pitstop authorisation: %w", err)
		}
		auths = append(auths, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return auths, nil
}

func (r *PitstopRepository) InsertAuthorisations(ctx context.Context, auths []*domain.PitstopAuthorisation) error {
	if len(auths) == 0 {
		return nil
	}

	// Build a bulk insert query
	valueStrings := make([]string, 0, len(auths))
	valueArgs := make([]interface{}, 0, len(auths)*10)

	for _, a := range auths {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			a.PitstopAuthID, a.DatasetID, a.DatasetName, a.UserID,
			a.RegulatorID, a.RegulatorName, a.MainconID, a.MainconName, a.Status, a.LastSyncedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO pitstop_authorisations (
			pitstop_auth_id, dataset_id, dataset_name, user_id, 
			regulator_id, regulator_name, maincon_id, maincon_name, status, last_synced_at
		) VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to insert pitstop authorisations: %w", err)
	}

	return nil
}

func (r *PitstopRepository) UpdateAuthorisations(ctx context.Context, auths []*domain.PitstopAuthorisation) error {
	if len(auths) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE pitstop_authorisations SET
			dataset_name = ?, user_id = ?, regulator_name = ?, maincon_name = ?,
			status = ?, last_synced_at = ?
		WHERE pitstop_auth_id = ?
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	for _, a := range auths {
		if _, err := stmt.ExecContext(ctx,
			a.DatasetName, a.UserID, a.RegulatorName, a.MainconName,
			a.Status, a.LastSyncedAt, a.PitstopAuthID,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update pitstop_auth_id %s: %w", a.PitstopAuthID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit update transaction: %w", err)
	}
	return nil
}
