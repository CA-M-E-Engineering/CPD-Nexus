package services

import (
	"context"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      ports.UserRepository
	analytics ports.AnalyticsService
}

func NewUserService(repo ports.UserRepository, analytics ports.AnalyticsService) ports.UserService {
	return &UserService{
		repo:      repo,
		analytics: analytics,
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User, password string) error {
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hash)
	}

	err := s.repo.Create(ctx, user)
	if err == nil {
		s.analytics.LogActivity(ctx, user.ID, "User Registered", "user", user.ID, fmt.Sprintf("New user account created for %s", user.Name))
	}
	return err
}

func (s *UserService) UpdateUser(ctx context.Context, id string, payload map[string]interface{}) error {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Manual patching from map - matching frontend keys and JSON tags
	if name, ok := payload["user_name"].(string); ok { user.Name = name }
	if username, ok := payload["username"].(string); ok { user.Username = username }
	if uType, ok := payload["user_type"].(string); ok { user.UserType = uType }
	if email, ok := payload["email"].(string); ok { user.ContactEmail = email }
	if phone, ok := payload["phone"].(string); ok { user.ContactPhone = phone }
	if status, ok := payload["status"].(string); ok { user.Status = status }
	if addr, ok := payload["address"].(string); ok { user.Address = addr }
	
	// Handle coordinates
	if lat, ok := payload["lat"].(float64); ok { user.Latitude = lat }
	if lng, ok := payload["lng"].(float64); ok { user.Longitude = lng }

	if bridgeWS, ok := payload["bridge_ws_url"].(string); ok { user.BridgeWSURL = &bridgeWS }
	if bridgeAuth, ok := payload["bridge_auth_token"].(string); ok { user.BridgeAuthToken = &bridgeAuth }
	if bridgeStat, ok := payload["bridge_status"].(string); ok { user.BridgeStatus = bridgeStat }

	if pwd, ok := payload["password"].(string); ok && pwd != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hash)
	}

	err = s.repo.Update(ctx, user)
	if err == nil {
		s.analytics.LogActivity(ctx, id, "User Updated", "user", id, "User profile and/or bridge configuration modified")
	}
	return err
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		s.analytics.LogActivity(ctx, id, "User Deleted", "user", id, "User account deactivated/removed")
	}
	return err
}
