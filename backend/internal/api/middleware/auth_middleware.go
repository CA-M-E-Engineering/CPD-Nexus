package middleware

import (
	"context"
	"net/http"
	"strings"

	"cpd-nexus/internal/core/ports"

	"github.com/golang-jwt/jwt/v5"
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
		var tokenStr string
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				tokenStr = cookie.Value
			}
		}

		if tokenStr == "" {
			// No token — allow request to continue but with no user scope.
			// RequireUserScope middleware will reject it if the route needs auth.
			next.ServeHTTP(w, r)
			return
		}

		claims, err := validateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized: invalid or expired token", http.StatusUnauthorized)
			return
		}

		userID, _ := claims["user_id"].(string)
		username, _ := claims["username"].(string)
		userType, _ := claims["user_type"].(string)

		ctx := context.WithValue(r.Context(), ports.UserIDKey, userID)
		ctx = context.WithValue(ctx, ports.UsernameKey, username)

		// Capture IP Address
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = strings.Split(r.RemoteAddr, ":")[0]
		} else {
			ip = strings.Split(ip, ",")[0]
		}
		ctx = context.WithValue(ctx, ports.IPAddressKey, strings.TrimSpace(ip))

		// Vendor role is derived from JWT claims. Only 'vendor' type has global system-wide visibility.
		// Standard client 'admin' or 'manager' roles are restricted to their own organization.
		if userType == "vendor" {
			ctx = context.WithValue(ctx, ports.IsVendorKey, true)
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

// RequireUserScope checks if userID is present in context, returns 401 if missing.
func RequireUserScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := ports.GetUserID(r.Context())
		if userID == "" {
			http.Error(w, "Unauthorized: valid JWT token required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdminScope checks if the user has admin/vendor privileges.
func RequireAdminScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ports.IsVendor(r.Context()) {
			http.Error(w, "Forbidden: administrative privileges required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
