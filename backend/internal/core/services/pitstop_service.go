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

	// 2. Load existing to maintain consistent ID for Upsert
	existingAuths, _ := s.pitstopRepo.GetAuthorisations(ctx)
	lookup := make(map[string]string)
	for _, e := range existingAuths {
		key := fmt.Sprintf("%s|%s|%s", e.DatasetID, e.RegulatorID, e.MainconID)
		lookup[key] = e.PitstopAuthID
	}

	var authorisations []*domain.PitstopAuthorisation
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
				authID, exists := lookup[key]
				if !exists {
					authID = fmt.Sprintf("pa%s%04d", idTimestamp, seq)
					seq++
				}

				// Create the domain entity
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
				authorisations = append(authorisations, auth)
			}
		}
	}

	// 3. Store into the database using Repository
	if len(authorisations) > 0 {
		return s.pitstopRepo.UpsertAuthorisations(ctx, authorisations)
	}

	return nil
}
