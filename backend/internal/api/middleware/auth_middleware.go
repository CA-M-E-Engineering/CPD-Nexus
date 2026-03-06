package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	IsVendorKey contextKey = "isVendor"
	UsernameKey contextKey = "username"
)

// JWTSecret is the shared secret for validating tokens.
// It must be set at startup via SetJWTSecret before any requests are handled.
var jwtSecret []byte

// SetJWTSecret configures the JWT signing secret used by the middleware.
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// UserScopeMiddleware validates the JWT from the Authorization header and populates
// the request context with user_id and isVendor status extracted from token claims.
// The X-User-ID header is intentionally ignored to prevent spoofing.
func UserScopeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// No token — allow request to continue but with no user scope.
			// RequireUserScope middleware will reject it if the route needs auth.
			next.ServeHTTP(w, r)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized: invalid or expired token", http.StatusUnauthorized)
			return
		}

		userID, _ := claims["user_id"].(string)
		username, _ := claims["username"].(string)
		userType, _ := claims["user_type"].(string)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UsernameKey, username)

		// Vendor/admin role is derived from JWT claims, not a hardcoded ID
		if userType == "admin" || userType == "vendor" {
			ctx = context.WithValue(ctx, IsVendorKey, true)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateJWT parses and validates a JWT token string, returning its claims.
func validateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

// GetUserID retrieves the userID from the context.
func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDKey).(string); ok {
		return v
	}
	return ""
}

// IsVendor checks if the current context belongs to a system vendor or admin.
func IsVendor(ctx context.Context) bool {
	if v, ok := ctx.Value(IsVendorKey).(bool); ok {
		return v
	}
	return false
}

// RequireUserScope checks if userID is present in context, returns 401 if missing.
func RequireUserScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		if userID == "" {
			http.Error(w, "Unauthorized: valid JWT token required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
