package domain

const (
	// Worker/Site/Project Status
	StatusActive   = "active"
	StatusInactive = "inactive"

	// Sync status values are defined in worker.go as SyncStatusPendingUpdate, SyncStatusSynced, SyncStatusPendingRegistration

	// User Types
	UserTypeUser  = "user"
	UserTypeAdmin = "admin"

	// Roles
	RoleWorker  = "worker"
	RolePIC     = "pic"
	RoleManager = "manager"
)
