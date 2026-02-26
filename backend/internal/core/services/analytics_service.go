package services

import (
	"context"
	"sgbuildex/internal/core/ports"
)

type AnalyticsService struct {
	repo ports.AnalyticsRepository
}

func NewAnalyticsService(repo ports.AnalyticsRepository) ports.AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) GetDashboardStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	return s.repo.GetDashboardStats(ctx, userID)
}

func (s *AnalyticsService) GetActivityLog(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	// Mock activity log for now
	logs := []map[string]interface{}{
		{"id": 1, "user": "System", "action": "User Dashboard Loaded", "target": userID, "time": "Just now"},
		{"id": 2, "user": "System", "action": "Daily Sync Check", "target": "Cloud", "time": "5 mins ago"},
	}
	return logs, nil
}

func (s *AnalyticsService) GetDetailedAnalytics(ctx context.Context, userID string) (map[string]interface{}, error) {
	return s.repo.GetDetailedAnalytics(ctx, userID)
}
