package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"

// UserScopeMiddleware extracts the user ID from the X-User-ID header OR query parameter
// and injects it into the request context.
func UserScopeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = r.URL.Query().Get("user_id")
		}

		if userID != "" {
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
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

// RequireUserScope checks if userID is present in context, returns 401 if missing
func RequireUserScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		if userID == "" {
			http.Error(w, "User-ID scope required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
