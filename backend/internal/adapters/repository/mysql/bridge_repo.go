package mysql

import (
	"context"
	"database/sql"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

type BridgeRepository struct {
	db *sql.DB
}

func NewBridgeRepository(db *sql.DB) ports.BridgeRepository {
	return &BridgeRepository{db: db}
}

func (r *BridgeRepository) GetActiveBridgeWorkers(ctx context.Context) ([]ports.BridgeWorkerTask, error) {
	query := `
		SELECT w.worker_id, w.user_id, p.site_id 
		FROM workers w
		JOIN projects p ON w.current_project_id = p.project_id
		WHERE w.status = ? AND w.current_project_id IS NOT NULL`

	rows, err := r.db.QueryContext(ctx, query, domain.StatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []ports.BridgeWorkerTask
	for rows.Next() {
		var t ports.BridgeWorkerTask
		if err := rows.Scan(&t.WorkerID, &t.UserID, &t.SiteID); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *BridgeRepository) GetActiveDeviceSNsBySite(ctx context.Context, siteID string) ([]string, error) {
	query := "SELECT sn FROM devices WHERE site_id = ? AND status != ?"
	rows, err := r.db.QueryContext(ctx, query, siteID, domain.StatusInactive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sns []string
	for rows.Next() {
		var sn string
		if err := rows.Scan(&sn); err == nil {
			sns = append(sns, sn)
		}
	}
	return sns, nil
}

func (r *BridgeRepository) GetWorkerOwnerID(ctx context.Context, workerID string) (string, error) {
	var ownerID string
	err := r.db.QueryRowContext(ctx, "SELECT user_id FROM workers WHERE worker_id = ?", workerID).Scan(&ownerID)
	return ownerID, err
}

func (r *BridgeRepository) GetActiveBridges(ctx context.Context) ([]ports.BridgeConfig, error) {
	query := "SELECT user_id, bridge_ws_url, bridge_auth_token FROM users WHERE bridge_status = ? AND bridge_ws_url IS NOT NULL"
	rows, err := r.db.QueryContext(ctx, query, domain.StatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []ports.BridgeConfig
	for rows.Next() {
		var c ports.BridgeConfig
		var authToken sql.NullString
		if err := rows.Scan(&c.UserID, &c.WSURL, &authToken); err != nil {
			continue
		}
		if authToken.Valid {
			c.AuthToken = authToken.String
		}
		configs = append(configs, c)
	}
	return configs, nil
}
