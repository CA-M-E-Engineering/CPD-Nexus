package services

import (
	"context"
	"cpd-nexus/internal/core/ports"
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

func (s *AnalyticsService) GetActivityLog(ctx context.Context, userID string, filters map[string]interface{}) ([]map[string]interface{}, error) {
	return s.repo.GetActivityLog(ctx, userID, filters)
}

func (s *AnalyticsService) LogActivity(ctx context.Context, userID, action, targetType, targetID, details string) error {
	// Identify the actor performing the action from context
	actorID := ports.GetUserID(ctx)
	actorName := ports.GetUsername(ctx)

	// Determine log "owner" (who sees it in their feed). 
	// Default to the provided userID, but override with actorID if available 
	// so users primarily see their own actions (as requested).
	ownerID := userID
	if actorID != "" {
		ownerID = actorID
	}

	if actorName == "" {
		// Fallback: If it's a login action or self-action and we don't have a name in ctx,
		// we can try to use the userID as a label if identity is clear
		if action == "Login" {
			actorName = userID
		} else if userID == "system" {
			actorName = "System"
		} else {
			actorName = "Anonymous"
		}
	}

	activity := map[string]interface{}{
		"user_id":     ownerID,
		"user_name":   actorName,
		"action":      action,
		"target_type": targetType,
		"target_id":   targetID,
		"details":     details,
		"ip_address":  ports.GetIPAddress(ctx),
	}
	return s.repo.LogActivity(ctx, activity)
}

func (s *AnalyticsService) GetDetailedAnalytics(ctx context.Context, userID string) (map[string]interface{}, error) {
	return s.repo.GetDetailedAnalytics(ctx, userID)
}
