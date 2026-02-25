package services

import (
	"context"
	"fmt"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, u *domain.User, password string) error {
	// 1. Generate ID
	// Note: UserService needs access to db for idgen if we don't have a cleaner way.
	// For now, let's assume the ID is generated outside or we need a better idgen port.
	// Actually, the current idgen takes *sql.DB which violates hexagonal.
	// I'll skip auto-ID generation here or use a simpler one for now.
	if u.ID == "" {
		return fmt.Errorf("user ID is required (modular ID generation pending)")
	}

	// 2. Generate Username if missing
	if u.Username == "" {
		u.Username = s.generateUsername(u.Name, u.ID)
	}

	// 3. Hash Password
	if password == "" {
		password = "password123"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.PasswordHash = string(hash)

	if u.Status == "" {
		u.Status = domain.StatusActive
	}

	return s.repo.Create(ctx, u)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, payload map[string]interface{}) error {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("user not found")
	}

	// Overlay logic
	if v, ok := payload["user_name"].(string); ok {
		existing.Name = v
	}
	if v, ok := payload["user_type"].(string); ok {
		existing.UserType = v
	}
	if v, ok := payload["email"].(string); ok {
		existing.ContactEmail = v
	}
	if v, ok := payload["phone"].(string); ok {
		existing.ContactPhone = v
	}
	if v, ok := payload["username"].(string); ok {
		existing.Username = v
	}
	if v, ok := payload["status"].(string); ok {
		existing.Status = v
	}
	if v, ok := payload["address"].(string); ok {
		existing.Address = v
	}
	if v, ok := payload["lat"].(float64); ok {
		existing.Latitude = v
	}
	if v, ok := payload["lng"].(float64); ok {
		existing.Longitude = v
	}

	if v, ok := payload["password"].(string); ok && v != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		existing.PasswordHash = string(hash)
	}

	return s.repo.Update(ctx, existing)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) generateUsername(name, id string) string {
	baseName := ""
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			baseName += strings.ToLower(string(char))
		} else if char == ' ' {
			baseName += "_"
		}
	}
	if len(baseName) > 15 {
		baseName = baseName[:15]
	}
	suffix := ""
	if len(id) >= 4 {
		suffix = id[len(id)-4:]
	}
	return baseName + "_" + suffix
}
