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

func (s *AttendanceService) GetAttendance(ctx context.Context, id string) (*domain.Attendance, error) {
	return s.repo.Get(ctx, id)
}

func (s *AttendanceService) ListAttendance(ctx context.Context, userID, siteID, workerID, date string) ([]domain.Attendance, error) {
	return s.repo.List(ctx, userID, siteID, workerID, date)
}

func (s *AttendanceService) ProcessBridgeAttendance(ctx context.Context, deviceSN string, personID string, timeIn, timeOut string, rawPayload []byte) error {
	// 1. Resolve Device
	device, err := s.deviceRepo.GetBySN(ctx, deviceSN)
	if err != nil || device == nil {
		return fmt.Errorf("failed to resolve device SN %s: %w", deviceSN, err)
	}

	// 2. Resolve Worker
	worker, err := s.workerRepo.GetByFIN(ctx, personID)
	if err != nil {
		return fmt.Errorf("database error resolving worker NRIC %s: %w", personID, err)
	}
	if worker == nil {
		return fmt.Errorf("worker NRIC %s not found in the database", personID)
	}

	// 3. Parse Times
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

	var siteID string
	if device.SiteID != nil {
		siteID = *device.SiteID
	}

	// 4. Create Attendance Record
	attendance := &domain.Attendance{
		ID:              s.generateNextID(ctx),
		DeviceID:        device.ID,
		WorkerID:        worker.ID,
		SiteID:          siteID,
		UserID:          device.UserID,
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
