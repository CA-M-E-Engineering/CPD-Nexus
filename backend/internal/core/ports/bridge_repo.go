package ports

import "context"

type BridgeWorkerTask struct {
	WorkerID string
	UserID   string
	SiteID   string
}

type BridgeConfig struct {
	UserID    string
	WSURL     string
	AuthToken string
}

type BridgeRepository interface {
	GetActiveBridgeWorkers(ctx context.Context) ([]BridgeWorkerTask, error)
	GetActiveDeviceSNsBySite(ctx context.Context, siteID string) ([]string, error)
	GetWorkerOwnerID(ctx context.Context, workerID string) (string, error)
	GetActiveBridges(ctx context.Context) ([]BridgeConfig, error)
	LogBridgeInteraction(ctx context.Context, userID, action, requestID string, requestPayload, responsePayload []byte, statusCode int) error
}
