package services

import (
	"context"
	"fmt"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"time"
)

type AttendanceService struct {
	repo       ports.AttendanceRepository
	workerRepo ports.WorkerRepository
	deviceRepo ports.DeviceRepository
	analytics  ports.AnalyticsService
	// No in-process state for ID generation — delegated to the DB for scale safety.
}

func NewAttendanceService(repo ports.AttendanceRepository, workerRepo ports.WorkerRepository, deviceRepo ports.DeviceRepository, analytics ports.AnalyticsService) ports.AttendanceService {
	return &AttendanceService{
		repo:       repo,
		workerRepo: workerRepo,
		deviceRepo: deviceRepo,
		analytics:  analytics,
	}
}

func (s *AttendanceService) GetAttendance(ctx context.Context, userID, id string) (*domain.Attendance, error) {
	return s.repo.Get(ctx, userID, id)
}

func (s *AttendanceService) ListAttendance(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error) {
	return s.repo.List(ctx, userID, siteID, workerID, date)
}

func (s *AttendanceService) UpdateAttendance(ctx context.Context, userID, id string, timeIn, timeOut *time.Time) error {
	if id == "" {
		return fmt.Errorf("attendance ID is required")
	}

	err := s.repo.Update(ctx, userID, id, timeIn, timeOut)
	if err == nil {
		s.analytics.LogActivity(ctx, userID, "Attendance Updated", "attendance", id, fmt.Sprintf("Updated attendance times for %s", id))
	}
	return err
}

func (s *AttendanceService) ProcessBridgeAttendance(ctx context.Context, workerID string, timeIn, timeOut string, rawPayload []byte) error {
	// 1. Resolve Worker
	// We use the internal workerID provided by the bridge (which we sent in the request)
	worker, err := s.workerRepo.Get(ctx, "", workerID)
	if err != nil {
		return fmt.Errorf("database error resolving worker ID %s: %w", workerID, err)
	}
	if worker == nil {
		return fmt.Errorf("worker ID %s not found in the database", workerID)
	}

	// 2. Parse Times
	tIn, err := time.Parse(time.RFC3339, timeIn)
	if err != nil {
		// Try fallback format if RFC3339 fails
		tIn, err = time.Parse("2006-01-02T15:04:05", timeIn)
	}
	var tOutPtr *time.Time
	if timeOut != "" {
		tOut, err := time.Parse(time.RFC3339, timeOut)
		if err != nil {
			tOut, _ = time.Parse("2006-01-02T15:04:05", timeOut)
		}
		tOutPtr = &tOut
	}

	// 3. Generate attendance ID atomically at the DB layer (scale-safe)
	attendanceID, err := s.repo.GenerateNextID(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate attendance ID: %w", err)
	}

	// 4. Create Attendance Record
	attendance := &domain.Attendance{
		ID:              attendanceID,
		DeviceID:        "BRIDGE_AGGREGATED", // No single device ID for aggregated records
		WorkerID:        worker.ID,
		SiteID:          worker.SiteID,
		UserID:          worker.UserID,
		TimeIn:          &tIn,
		TimeOut:         tOutPtr,
		Direction:       "unknown",
		TradeCode:       worker.PersonTrade,
		Status:          "pending",
		SubmissionDate:  tIn.Format("2006-01-02"),
		ResponsePayload: string(rawPayload),
	}

	err = s.repo.Create(ctx, attendance)
	if err == nil {
		s.analytics.LogActivity(ctx, worker.UserID, "Attendance Logged", "worker", worker.ID, fmt.Sprintf("Aggregated attendance processed for %s at site %s", worker.Name, worker.SiteID))
	}
	return err
}
