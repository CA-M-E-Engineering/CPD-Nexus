package services

import (
	"context"
	"testing"
	"time"

	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// ─────────────────────────────────────────────
// Mocks matching actual port interfaces
// ─────────────────────────────────────────────

type authTestUserRepo struct {
	mock.Mock
}

func (m *authTestUserRepo) Get(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *authTestUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *authTestUserRepo) List(ctx context.Context) ([]domain.User, error)           { return nil, nil }
func (m *authTestUserRepo) Create(ctx context.Context, u *domain.User) error           { return nil }
func (m *authTestUserRepo) Update(ctx context.Context, u *domain.User) error           { return nil }
func (m *authTestUserRepo) Delete(ctx context.Context, id string) error                { return nil }

// analytics matching ports.AnalyticsService exactly
type authTestAnalytics struct {
	mock.Mock
}

func (m *authTestAnalytics) LogActivity(ctx context.Context, userID, action, targetType, targetID, details string) error {
	args := m.Called(ctx, userID, action, targetType, targetID, details)
	return args.Error(0)
}

func (m *authTestAnalytics) GetDashboardStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	return nil, nil
}

func (m *authTestAnalytics) GetActivityLog(ctx context.Context, userID string, filters map[string]interface{}) ([]map[string]interface{}, error) {
	return nil, nil
}

func (m *authTestAnalytics) GetDetailedAnalytics(ctx context.Context, userID string) (map[string]interface{}, error) {
	return nil, nil
}

func (m *authTestAnalytics) SetUserRepo(repo ports.UserRepository) {}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func authTestHashPwd(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	return string(hash)
}

func authTestParseJWT(t *testing.T, tokenStr, secret string) jwt.MapClaims {
	t.Helper()
	token, err := jwt.Parse(tokenStr, func(tk *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !token.Valid {
		t.Fatalf("parse JWT: %v", err)
	}
	return token.Claims.(jwt.MapClaims)
}

// ─────────────────────────────────────────────
// Tests
// ─────────────────────────────────────────────

func TestAuthService_Login_Success(t *testing.T) {
	repo := &authTestUserRepo{}
	analytics := &authTestAnalytics{}

	u := &domain.User{ID: "u-001", Username: "admin", UserType: "vendor", Status: "active"}
	u.PasswordHash = authTestHashPwd(t, "testpass")
	repo.On("GetByUsername", mock.Anything, "admin").Return(u, nil)
	analytics.On("LogActivity", mock.Anything, "u-001", "Login", "user", "u-001", mock.Anything).Return(nil)

	svc := NewAuthService(repo, "secret", analytics)
	token, user, err := svc.Login(context.Background(), "admin", "testpass")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "admin", user.Username)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	repo := &authTestUserRepo{}
	analytics := &authTestAnalytics{}

	u := &domain.User{ID: "u-001", Username: "admin", Status: "active"}
	u.PasswordHash = authTestHashPwd(t, "correct")
	repo.On("GetByUsername", mock.Anything, "admin").Return(u, nil)

	svc := NewAuthService(repo, "secret", analytics)
	_, _, err := svc.Login(context.Background(), "admin", "wrong")

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := &authTestUserRepo{}
	analytics := &authTestAnalytics{}
	repo.On("GetByUsername", mock.Anything, "nobody").Return(nil, nil)

	svc := NewAuthService(repo, "secret", analytics)
	_, _, err := svc.Login(context.Background(), "nobody", "x")

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestAuthService_Login_JWTClaimsAndExpiry(t *testing.T) {
	repo := &authTestUserRepo{}
	analytics := &authTestAnalytics{}

	u := &domain.User{ID: "u-abc", Username: "testuser", UserType: "client", Status: "active"}
	u.PasswordHash = authTestHashPwd(t, "mypass")
	repo.On("GetByUsername", mock.Anything, "testuser").Return(u, nil)
	analytics.On("LogActivity", mock.Anything, mock.Anything, "Login", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	svc := NewAuthService(repo, "my-secret", analytics)
	token, _, err := svc.Login(context.Background(), "testuser", "mypass")
	assert.NoError(t, err)

	claims := authTestParseJWT(t, token, "my-secret")

	assert.Equal(t, "u-abc", claims["user_id"])
	assert.Equal(t, "client", claims["user_type"])
	assert.Equal(t, "testuser", claims["username"])

	// Token should expire ~2 hours from now (tolerance: 1h–3h)
	exp := int64(claims["exp"].(float64))
	expTime := time.Unix(exp, 0)
	assert.True(t, expTime.After(time.Now().Add(1*time.Hour)), "token should not expire in under 1h")
	assert.True(t, expTime.Before(time.Now().Add(3*time.Hour)), "token should not live beyond 3h")
}

func TestAuthService_Login_NoAuditLogOnFailure(t *testing.T) {
	repo := &authTestUserRepo{}
	analytics := &authTestAnalytics{}

	u := &domain.User{ID: "u-001", Username: "admin", Status: "active"}
	u.PasswordHash = authTestHashPwd(t, "correct")
	repo.On("GetByUsername", mock.Anything, "admin").Return(u, nil)

	svc := NewAuthService(repo, "secret", analytics)
	_, _, err := svc.Login(context.Background(), "admin", "wrong")
	assert.Error(t, err)

	// LogActivity must NOT be called when authentication fails
	analytics.AssertNotCalled(t, "LogActivity")
}
