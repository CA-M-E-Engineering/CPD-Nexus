package services

import (
	"context"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo            ports.UserRepository
	analytics       ports.AnalyticsService
	defaultPassword string
}

func NewUserService(repo ports.UserRepository, analytics ports.AnalyticsService, defaultPassword string) ports.UserService {
	return &UserService{
		repo:            repo,
		analytics:       analytics,
		defaultPassword: defaultPassword,
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User, password string) error {
	// If no password provided, use the global default password to avoid hardcoding specific user credentials
	finalPassword := password
	if finalPassword == "" {
		finalPassword = s.defaultPassword
	}

	if finalPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(finalPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hash)
	}
    
    // Auto-generate Bridge Config if missing
    if user.BridgeAuthToken == nil || *user.BridgeAuthToken == "" {
        token := generateSecureToken(16)
        user.BridgeAuthToken = &token
    }
    
    // Default bridge status to active for new users to simplify setup
    if user.BridgeStatus == "" {
        user.BridgeStatus = "active"
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

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
