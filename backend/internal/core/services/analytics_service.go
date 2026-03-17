package services

import (
	"context"
	"cpd-nexus/internal/core/ports"
)

type AnalyticsService struct {
	repo     ports.AnalyticsRepository
	userRepo ports.UserRepository
}

func NewAnalyticsService(repo ports.AnalyticsRepository) ports.AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) SetUserRepo(repo ports.UserRepository) {
	s.userRepo = repo
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
	ownerID := userID
	if actorID != "" {
		ownerID = actorID
	}

	// PROACTIVE ENHANCEMENT: Resolve names if they are IDs or missing
	if (actorName == "" || actorName == actorID) && actorID != "" && s.userRepo != nil {
		if u, err := s.userRepo.Get(ctx, actorID); err == nil && u != nil {
			actorName = u.Name
		}
	}

	// FALLBACKS
	if actorName == "" {
		if action == "Login" && userID != "" && s.userRepo != nil {
			// For login, the actorID is not yet in context, use passed userID
			if u, err := s.userRepo.Get(ctx, userID); err == nil && u != nil {
				actorName = u.Name
			} else {
				actorName = userID
			}
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
