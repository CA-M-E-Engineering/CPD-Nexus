package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type AttendanceRepository interface {
	Get(ctx context.Context, id string) (*domain.Attendance, error)
	List(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error)
	Create(ctx context.Context, a *domain.Attendance) error
	GetMaxID(ctx context.Context, pattern string) (string, error)
}

type AttendanceService interface {
	GetAttendance(ctx context.Context, id string) (*domain.Attendance, error)
	ListAttendance(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error)
	ProcessBridgeAttendance(ctx context.Context, deviceSN string, personID string, timeIn, timeOut string, rawPayload []byte) error
}
