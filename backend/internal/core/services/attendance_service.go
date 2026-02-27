package services

import (
	"context"
	"fmt"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"strings"
	"sync"
	"time"
)

type AttendanceService struct {
	repo       ports.AttendanceRepository
	workerRepo ports.WorkerRepository
	deviceRepo ports.DeviceRepository

	mu      sync.Mutex
	lastSeq int
	lastDay string
}

func NewAttendanceService(repo ports.AttendanceRepository, workerRepo ports.WorkerRepository, deviceRepo ports.DeviceRepository) ports.AttendanceService {
	return &AttendanceService{
		repo:       repo,
		workerRepo: workerRepo,
		deviceRepo: deviceRepo,
	}
}

func (s *AttendanceService) GetAttendance(ctx context.Context, userID, id string) (*domain.Attendance, error) {
	return s.repo.Get(ctx, userID, id)
}

func (s *AttendanceService) ListAttendance(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error) {
	return s.repo.List(ctx, userID, siteID, workerID, date)
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

	// 3. Create Attendance Record
	attendance := &domain.Attendance{
		ID:              s.generateNextID(ctx),
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

	return s.repo.Create(ctx, attendance)
}

func (s *AttendanceService) generateNextID(ctx context.Context) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	day := time.Now().Format("20060102")

	if s.lastDay != day {
		s.lastDay = day
		s.lastSeq = 0

		maxID, err := s.repo.GetMaxID(ctx, "ATT-"+day+"-%")
		if err == nil && maxID != "" {
			parts := strings.Split(maxID, "-")
			if len(parts) == 3 {
				var seq int
				fmt.Sscanf(parts[2], "%d", &seq)
				s.lastSeq = seq
			}
		}
	}

	s.lastSeq++
	return fmt.Sprintf("ATT-%s-%04d", day, s.lastSeq)
}
