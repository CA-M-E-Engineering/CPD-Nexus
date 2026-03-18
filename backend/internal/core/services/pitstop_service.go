package services

import (
	"context"
	"fmt"

	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"time"
)

// PitstopService orchestrates pitstop authorisation management and attendance submission.
// It depends only on port interfaces — never on concrete adapter types.
type PitstopService struct {
	pitstopRepo    ports.PitstopRepository
	externalClient ports.ExternalSubmitter // was *sgbuildex.Client — now decoupled via interface (#5)
	attendanceRepo ports.AttendanceRepository
	submissionRepo ports.SubmissionRepository
	settingsRepo   ports.SettingsRepository
	analytics      ports.AnalyticsService
}

func NewPitstopService(
	repo ports.PitstopRepository,
	client ports.ExternalSubmitter,
	attendanceRepo ports.AttendanceRepository,
	submissionRepo ports.SubmissionRepository,
	settingsRepo ports.SettingsRepository,
	analytics ports.AnalyticsService,
) *PitstopService {
	return &PitstopService{
		pitstopRepo:    repo,
		externalClient: client,
		attendanceRepo: attendanceRepo,
		submissionRepo: submissionRepo,
		settingsRepo:   settingsRepo,
		analytics:      analytics,
	}
}

// GetAuthorisations returns the currently stored authorisations. (Globally visible for Vendors)
func (s *PitstopService) GetAuthorisations(ctx context.Context, userID string) ([]*domain.PitstopAuthorisation, error) {
	return s.pitstopRepo.GetAuthorisations(ctx, userID)
}

// SyncConfig fetches the newest configs from the Pitstop API and upserts them
func (s *PitstopService) SyncConfig(ctx context.Context, userID string) error {
	// 1. Fetch from Pitstop API via the port interface — no concrete adapter type referenced
	cfgResponse, err := s.externalClient.FetchPitstopConfig(ctx)
	if err != nil {
		return fmt.Errorf("pitstop API fetch failed: %w", err)
	}

	// 2. Load existing to maintain consistent ID for Insert/Update checks
	existingAuths, _ := s.pitstopRepo.GetAuthorisations(ctx, "")
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

	// 3. Map response to domain entities — using ports-level response types
	for _, produce := range cfgResponse.Produces {
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
					pitstopAuthID := fmt.Sprintf("pa%s%04d", idTimestamp, seq)
					seq++
					toInsert = append(toInsert, &domain.PitstopAuthorisation{
						PitstopAuthID:  pitstopAuthID,
						DatasetID:      datasetID,
						DatasetName:    datasetName,
						RegulatorID:    regulatorID,
						RegulatorName:  regulatorName,
						OnBehalfOfID:   onBehalfOfID,
						OnBehalfOfName: onBehalfOfName,
						Status:         "ACTIVE",
						LastSyncedAt:   &now,
						UserID:         &userID,
					})
				}
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

	s.analytics.LogActivity(ctx, userID, "Pitstop Sync", "system", userID, "Synchronized organizational authorisations with SGBuildex Pitstop API")
	return nil
}

// GetProjectsWithPendingAttendance returns a list of unique projects that have pending attendance records
func (s *PitstopService) GetProjectsWithPendingAttendance(ctx context.Context, userID string) ([]domain.Project, error) {
	return s.attendanceRepo.ExtractProjectsWithPendingAttendance(ctx, userID)
}

// TestSubmission extracts pending attendance for a given project and immediately pushes it
func (s *PitstopService) TestSubmission(ctx context.Context, userID, projectID string) (submittedCount int, failedCount int, err error) {
	settings, err := s.loadSettings(ctx)
	if err != nil {
		return 0, 0, err
	}

	rows, err := s.attendanceRepo.ExtractPendingAttendanceByProject(ctx, userID, projectID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to extract project attendance: %w", err)
	}

	if len(rows) == 0 {
		return 0, 0, nil
	}

	submittedCount, failedCount, err = s.externalClient.SubmitManpowerUtilization(ctx, s.submissionRepo, settings, rows)

	details := fmt.Sprintf("Submitted %d records (%d validation failed) for project %s", submittedCount, failedCount, projectID)
	if err != nil {
		details = fmt.Sprintf("Submission failed: %v. Records processed: %d", err, submittedCount)
	}

	s.analytics.LogActivity(ctx, userID, "Manual CPD Submission", "project", projectID, details)

	if err != nil {
		return submittedCount, failedCount, fmt.Errorf("failed to submit payloads: %w", err)
	}
	return submittedCount, failedCount, nil
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

	// Submit via the port interface — no concrete adapter type referenced
	_, _, err = s.externalClient.SubmitManpowerUtilization(ctx, s.submissionRepo, settings, rows)
	if err == nil {
		s.analytics.LogActivity(ctx, "system", "Scheduled CPD Submission", "system", "pitstop", fmt.Sprintf("Automatically submitted %d pending attendance records to SGBuildex", len(rows)))
	}
	return err
}

// AssignOnBehalfOfToUser assigns a set of contractor names to a specific UserID for pitstop authorisations
func (s *PitstopService) AssignOnBehalfOfToUser(ctx context.Context, userID string, onBehalfOfNames []string) error {
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	err := s.pitstopRepo.AssignOnBehalfOfToUser(ctx, userID, onBehalfOfNames)
	if err == nil {
		actorUserID := ports.GetUserID(ctx)
		s.analytics.LogActivity(ctx, actorUserID, "Pitstop Assignment", "user", userID, fmt.Sprintf("Assigned %d contractor authorisations to user %s", len(onBehalfOfNames), userID))
	}
	return err
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
