package domain

const (
	// Worker/Site/Project Status
	StatusActive   = "active"
	StatusInactive = "inactive"

	// Sync Status
	SyncStatusReady   = 1
	SyncStatusPending = 2
	SyncStatusUpdate  = 0

	// User Types
	UserTypeUser  = "user"
	UserTypeAdmin = "admin"

	// Roles
	RoleWorker  = "worker"
	RolePIC     = "pic"
	RoleManager = "manager"
)
