package ports

import (
	"context"
)

type AnalyticsRepository interface {
	GetDashboardStats(ctx context.Context, userID string) (map[string]interface{}, error)
	GetDetailedAnalytics(ctx context.Context, userID string) (map[string]interface{}, error)
	LogActivity(ctx context.Context, log map[string]interface{}) error
	GetActivityLog(ctx context.Context, userID string, filters map[string]interface{}) ([]map[string]interface{}, error)
}

type AnalyticsService interface {
	GetDashboardStats(ctx context.Context, userID string) (map[string]interface{}, error)
	GetActivityLog(ctx context.Context, userID string, filters map[string]interface{}) ([]map[string]interface{}, error)
	GetDetailedAnalytics(ctx context.Context, userID string) (map[string]interface{}, error)
	LogActivity(ctx context.Context, userID, action, targetType, targetID, details string) error
}
