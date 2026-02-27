package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

type AttendanceRepository interface {
	Get(ctx context.Context, userID, id string) (*domain.Attendance, error)
	List(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error)
	Create(ctx context.Context, a *domain.Attendance) error
	GetMaxID(ctx context.Context, pattern string) (string, error)
	ExtractPendingAttendance(ctx context.Context) ([]domain.AttendanceRow, error)
	ExtractMonthlyDistributionData(ctx context.Context) ([]domain.MonthlyDistributionRow, error)
}

type AttendanceService interface {
	GetAttendance(ctx context.Context, userID, id string) (*domain.Attendance, error)
	ListAttendance(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error)
	ProcessBridgeAttendance(ctx context.Context, workerID string, timeIn, timeOut string, rawPayload []byte) error
}
