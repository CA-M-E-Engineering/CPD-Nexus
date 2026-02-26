package ports

import (
	"context"
)

type AnalyticsRepository interface {
	GetDashboardStats(ctx context.Context, userID string) (map[string]interface{}, error)
	GetDetailedAnalytics(ctx context.Context, userID string) (map[string]interface{}, error)
}

type AnalyticsService interface {
	GetDashboardStats(ctx context.Context, userID string) (map[string]interface{}, error)
	GetActivityLog(ctx context.Context, userID string) ([]map[string]interface{}, error)
	GetDetailedAnalytics(ctx context.Context, userID string) (map[string]interface{}, error)
}
