package ports

import (
	"context"
	"cpd-nexus/internal/core/domain"
	"time"
)

type AttendanceRepository interface {
	Get(ctx context.Context, userID, id string) (*domain.Attendance, error)
	List(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error)
	Create(ctx context.Context, a *domain.Attendance) error
	GetMaxID(ctx context.Context, pattern string) (string, error)
	// GenerateNextID atomically generates the next sequential attendance ID for a given day.
	// This is implemented at the DB layer to be safe for multi-instance deployments.
	GenerateNextID(ctx context.Context) (string, error)
	Update(ctx context.Context, userID, id string, timeIn, timeOut *time.Time) error
	ExtractPendingAttendance(ctx context.Context) ([]domain.AttendanceRow, error)
	ExtractPendingAttendanceByProject(ctx context.Context, userID, projectID string) ([]domain.AttendanceRow, error)
	ExtractProjectsWithPendingAttendance(ctx context.Context, userID string) ([]domain.Project, error)
}

type AttendanceService interface {
	GetAttendance(ctx context.Context, userID, id string) (*domain.Attendance, error)
	ListAttendance(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error)
	ProcessBridgeAttendance(ctx context.Context, workerID string, timeIn, timeOut string, rawPayload []byte) error
	UpdateAttendance(ctx context.Context, userID, id string, timeIn, timeOut *time.Time) error
}
