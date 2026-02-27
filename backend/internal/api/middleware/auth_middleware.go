package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	IsVendorKey contextKey = "isVendor"
)

// UserScopeMiddleware extracts the user ID and checks if it's the system vendor.
func UserScopeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = r.URL.Query().Get("user_id")
		}

		if userID != "" {
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			// Hardcoded God-Mode check for the primary vendor account
			if userID == "Owner_001" {
				ctx = context.WithValue(ctx, IsVendorKey, true)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserID retrieves the userID from the context
func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDKey).(string); ok {
		return v
	}
	return ""
}

// IsVendor checks if the current context belongs to a system vendor
func IsVendor(ctx context.Context) bool {
	if v, ok := ctx.Value(IsVendorKey).(bool); ok {
		return v
	}
	return false
}

// RequireUserScope checks if userID is present in context, returns 401 if missing
func RequireUserScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		if userID == "" {
			// Special bypass for vendor logged in checks if we had a proper session/JWT
			// For now, let's keep it strict but add debugging info to the response
			http.Error(w, "User-ID scope required (X-User-ID header missing or invalid)", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
