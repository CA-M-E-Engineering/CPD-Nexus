package services

import (
	"context"
	"fmt"
	"sgbuildex/internal/adapters/external/sgbuildex"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"time"
)

type PitstopService struct {
	pitstopRepo    ports.PitstopRepository
	pitstopClient  *sgbuildex.Client
	attendanceRepo ports.AttendanceRepository
	submissionRepo ports.SubmissionRepository
	settingsRepo   ports.SettingsRepository
}

func NewPitstopService(
	repo ports.PitstopRepository,
	client *sgbuildex.Client,
	attendanceRepo ports.AttendanceRepository,
	submissionRepo ports.SubmissionRepository,
	settingsRepo ports.SettingsRepository,
) *PitstopService {
	return &PitstopService{
		pitstopRepo:    repo,
		pitstopClient:  client,
		attendanceRepo: attendanceRepo,
		submissionRepo: submissionRepo,
		settingsRepo:   settingsRepo,
	}
}

// GetAuthorisations returns the currently stored authorisations. (Globally visible for Vendors)
func (s *PitstopService) GetAuthorisations(ctx context.Context) ([]*domain.PitstopAuthorisation, error) {
	return s.pitstopRepo.GetAuthorisations(ctx)
}

// SyncConfig fetches the newest configs from the Pitstop API and upserts them
func (s *PitstopService) SyncConfig(ctx context.Context, userID string) error {
	// 1. Fetch from Pitstop API
	cfgResponse, err := s.pitstopClient.FetchConfig(ctx)
	if err != nil {
		return fmt.Errorf("pitstop API fetch failed: %w", err)
	}

	// 2. Load existing to maintain consistent ID for Insert/Update checks
	existingAuths, _ := s.pitstopRepo.GetAuthorisations(ctx)
	existingMap := make(map[string]*domain.PitstopAuthorisation)
	for _, e := range existingAuths {
		key := fmt.Sprintf("%s|%s|%s", e.DatasetID, e.RegulatorID, e.OnBehalfOfID)
		existingMap[key] = e
	}

	var toInsert []*domain.PitstopAuthorisation
	var toUpdate []*domain.PitstopAuthorisation
	now := time.Now()
	idTimestamp := now.Format("20060102150405")
	seq := 1

	seenKeys := make(map[string]bool)

	// 3. Map JSON Response to Domain Entities
	for _, produce := range cfgResponse.Data.Produces {
		datasetID := produce.ID
		datasetName := produce.Name

		for _, to := range produce.To {
			regulatorID := to.ID
			regulatorName := to.Name

			for _, behalf := range to.OnBehalfOf {
				onBehalfOfID := behalf.ID
				onBehalfOfName := behalf.Name

				key := fmt.Sprintf("%s|%s|%s", datasetID, regulatorID, onBehalfOfID)
				seenKeys[key] = true

				if existing, exists := existingMap[key]; exists {
					modified := false
					if existing.DatasetName != datasetName ||
						existing.RegulatorName != regulatorName ||
						existing.OnBehalfOfName != onBehalfOfName ||
						existing.Status != "ACTIVE" {
						modified = true
					}

					if existing.UserID == nil || *existing.UserID == "" {
						existing.UserID = &userID
						modified = true
					}

					if modified {
						existing.DatasetName = datasetName
						existing.RegulatorName = regulatorName
						existing.OnBehalfOfName = onBehalfOfName
						existing.Status = "ACTIVE"
						existing.LastSyncedAt = &now
						toUpdate = append(toUpdate, existing)
					}
				} else {
					authID := fmt.Sprintf("pa%s%04d", idTimestamp, seq)
					seq++

					auth := &domain.PitstopAuthorisation{
						PitstopAuthID:  authID,
						DatasetID:      datasetID,
						DatasetName:    datasetName,
						UserID:         &userID,
						RegulatorID:    regulatorID,
						RegulatorName:  regulatorName,
						OnBehalfOfID:   onBehalfOfID,
						OnBehalfOfName: onBehalfOfName,
						Status:         "ACTIVE",
						LastSyncedAt:   &now,
					}
					toInsert = append(toInsert, auth)
				}
			}
		}
	}

	// 3.5 Mark entries no longer in Pitstop config as INACTIVE
	for key, existing := range existingMap {
		if !seenKeys[key] {
			if existing.Status != "INACTIVE" {
				existing.Status = "INACTIVE"
				existing.LastSyncedAt = &now
				toUpdate = append(toUpdate, existing)
			}
		}
	}

	// 4. Persist changes via Repository
	if len(toInsert) > 0 {
		if err := s.pitstopRepo.InsertAuthorisations(ctx, toInsert); err != nil {
			return err
		}
	}

	if len(toUpdate) > 0 {
		if err := s.pitstopRepo.UpdateAuthorisations(ctx, toUpdate); err != nil {
			return err
		}
	}

	return nil
}

// GetProjectsWithPendingAttendance returns a list of unique projects that have pending attendance records
func (s *PitstopService) GetProjectsWithPendingAttendance(ctx context.Context) ([]domain.Project, error) {
	return s.attendanceRepo.ExtractProjectsWithPendingAttendance(ctx)
}

// TestSubmission extracts pending attendance for a given project and immediately pushes it
func (s *PitstopService) TestSubmission(ctx context.Context, projectID string) (int, error) {
	settings, err := s.loadSettings(ctx)
	if err != nil {
		return 0, err
	}

	rows, err := s.attendanceRepo.ExtractPendingAttendanceByProject(ctx, projectID)
	if err != nil {
		return 0, fmt.Errorf("failed to extract project attendance: %w", err)
	}

	if len(rows) == 0 {
		return 0, nil
	}

	muSubmittables := s.mapRowsToSubmittables(rows)
	if err := sgbuildex.SubmitPayloads(ctx, s.submissionRepo, s.pitstopClient, settings, muSubmittables); err != nil {
		return 0, fmt.Errorf("failed to submit payloads: %w", err)
	}

	return len(muSubmittables), nil
}

// SubmitPendingAttendance extracts all non-submitted attendance and pushes it to SGBuildex.
// This is the method called by the scheduled task in main.go.
func (s *PitstopService) SubmitPendingAttendance(ctx context.Context) error {
	settings, err := s.loadSettings(ctx)
	if err != nil {
		return err
	}

	rows, err := s.attendanceRepo.ExtractPendingAttendance(ctx)
	if err != nil {
		return fmt.Errorf("failed to extract pending attendance: %w", err)
	}

	if len(rows) == 0 {
		return nil
	}

	muSubmittables := s.mapRowsToSubmittables(rows)
	return sgbuildex.SubmitPayloads(ctx, s.submissionRepo, s.pitstopClient, settings, muSubmittables)
}

// AssignOnBehalfOfToUser assigns a set of contractor names to a specific UserID for pitstop authorisations
func (s *PitstopService) AssignOnBehalfOfToUser(ctx context.Context, userID string, onBehalfOfNames []string) error {
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	return s.pitstopRepo.AssignOnBehalfOfToUser(ctx, userID, onBehalfOfNames)
}

// --- private helpers ---

// loadSettings fetches system settings, falling back to safe defaults on error.
func (s *PitstopService) loadSettings(ctx context.Context) (*domain.SystemSettings, error) {
	settings, err := s.settingsRepo.GetSettings(ctx)
	if err != nil {
		// Return defaults rather than failing — submission is best-effort
		return &domain.SystemSettings{
			MaxWorkersPerRequest: 100,
			MaxPayloadSizeKB:     256,
			MaxRequestsPerMinute: 150,
		}, nil
	}
	return settings, nil
}

// mapRowsToSubmittables converts attendance rows to typed Submittable wrappers.
func (s *PitstopService) mapRowsToSubmittables(rows []domain.AttendanceRow) []sgbuildex.ManpowerUtilizationWrapper {
	muPayloads := sgbuildex.MapAttendanceToManpower(rows)
	wrappers := make([]sgbuildex.ManpowerUtilizationWrapper, len(muPayloads))
	for i, p := range muPayloads {
		wrappers[i] = sgbuildex.ManpowerUtilizationWrapper{ManpowerUtilization: p}
	}
	return wrappers
}
