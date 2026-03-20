package mysql

import (
	"context"
	"database/sql"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/pkg/logger"
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
		JOIN users u ON w.user_id = u.user_id
		WHERE w.status = ? AND w.current_project_id IS NOT NULL AND u.bridge_status = ?`

	rows, err := r.db.QueryContext(ctx, query, domain.StatusActive, domain.StatusActive)
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

func (r *BridgeRepository) LogBridgeInteraction(ctx context.Context, userID, action, requestID string, requestPayload, responsePayload []byte, statusCode int) error {
	query := `
		INSERT INTO bridge_logs (user_id, action, request_id, request_payload, response_payload, status_code)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			response_payload = IF(VALUES(response_payload) IS NOT NULL, VALUES(response_payload), response_payload),
			status_code = IF(VALUES(status_code) IS NOT NULL, VALUES(status_code), status_code),
			updated_at = CURRENT_TIMESTAMP`

	var reqPl, respPl interface{}
	if len(requestPayload) > 0 {
		reqPl = requestPayload
	}
	if len(responsePayload) > 0 {
		respPl = responsePayload
	}

	var sc interface{}
	if statusCode != 0 {
		sc = statusCode
	}

	// 6 parameters total: user_id, action, request_id, request_payload, response_payload, status_code
	// The ON DUPLICATE KEY UPDATE uses VALUES() to refer to these.
	_, err := r.db.ExecContext(ctx, query, userID, action, requestID, reqPl, respPl, sc)
	if err != nil {
		logger.Errorf("BridgeRepo: Failed to log interaction (user: %s, action: %s, reqID: %s): %v", userID, action, requestID, err)
	}
	return err
}
