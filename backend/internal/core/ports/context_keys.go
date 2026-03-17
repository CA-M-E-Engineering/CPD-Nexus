package ports

import "context"

type ContextKey string

const (
	UserIDKey   ContextKey = "userID"
	IsVendorKey ContextKey = "isVendor"
	UsernameKey  ContextKey = "username"
	IPAddressKey ContextKey = "ipAddress"
)

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

// GetUsername retrieves the username from the context.
func GetUsername(ctx context.Context) string {
	if v, ok := ctx.Value(UsernameKey).(string); ok {
		return v
	}
	return ""
}

// GetIPAddress retrieves the client's IP address from the context.
func GetIPAddress(ctx context.Context) string {
	if v, ok := ctx.Value(IPAddressKey).(string); ok {
		return v
	}
	return ""
}
