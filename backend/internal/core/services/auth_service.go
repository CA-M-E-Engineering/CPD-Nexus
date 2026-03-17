package services

import (
	"context"
	"errors"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo             ports.UserRepository
	jwtSecret        string
	analyticsService ports.AnalyticsService
}

func NewAuthService(repo ports.UserRepository, jwtSecret string, analytics ports.AnalyticsService) ports.AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret, analyticsService: analytics}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Verify password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Issue a real signed JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"user_type": user.UserType,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	})

	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, errors.New("failed to issue token")
	}

	s.analyticsService.LogActivity(ctx, user.ID, "Login", "user", user.ID, "User logged in to the system")
	return tokenStr, user, nil
}
