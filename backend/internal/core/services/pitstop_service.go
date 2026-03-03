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
	pitstopRepo   ports.PitstopRepository
	pitstopClient *sgbuildex.Client
}

func NewPitstopService(repo ports.PitstopRepository, client *sgbuildex.Client) *PitstopService {
	return &PitstopService{
		pitstopRepo:   repo,
		pitstopClient: client,
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
		key := fmt.Sprintf("%s|%s|%s", e.DatasetID, e.RegulatorID, e.MainconID)
		existingMap[key] = e
	}

	var toInsert []*domain.PitstopAuthorisation
	var toUpdate []*domain.PitstopAuthorisation
	now := time.Now()
	idTimestamp := now.Format("20060102150405")
	seq := 1

	// 3. Map JSON Response to Domain Entities
	for _, produce := range cfgResponse.Data.Produces {
		datasetID := produce.ID
		datasetName := produce.Name

		for _, to := range produce.To {
			regulatorID := to.ID
			regulatorName := to.Name

			for _, behalf := range to.OnBehalfOf {
				mainconID := behalf.ID
				mainconName := behalf.Name

				key := fmt.Sprintf("%s|%s|%s", datasetID, regulatorID, mainconID)

				if existing, exists := existingMap[key]; exists {
					// Check if data actually changed
					modified := false
					if existing.DatasetName != datasetName ||
						existing.RegulatorName != regulatorName ||
						existing.MainconName != mainconName ||
						existing.Status != "ACTIVE" {
						modified = true
					}

					// Apply updates only if modified from the payload changes
					if modified {
						existing.DatasetName = datasetName
						existing.RegulatorName = regulatorName
						existing.MainconName = mainconName
						existing.Status = "ACTIVE"
						existing.LastSyncedAt = &now
						if existing.UserID == nil || *existing.UserID != userID {
							existing.UserID = &userID
						}
						toUpdate = append(toUpdate, existing)
					}
				} else {
					// Insert New
					authID := fmt.Sprintf("pa%s%04d", idTimestamp, seq)
					seq++

					auth := &domain.PitstopAuthorisation{
						PitstopAuthID: authID,
						DatasetID:     datasetID,
						DatasetName:   datasetName,
						UserID:        &userID,
						RegulatorID:   regulatorID,
						RegulatorName: regulatorName,
						MainconID:     mainconID,
						MainconName:   mainconName,
						Status:        "ACTIVE",
						LastSyncedAt:  &now,
					}
					toInsert = append(toInsert, auth)
				}
			}
		}
	}

	// 4. Store into the database using Repository
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
