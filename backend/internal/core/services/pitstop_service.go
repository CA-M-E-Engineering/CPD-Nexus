package services

import (
	"context"
	"fmt"
	"sgbuildex/internal/adapters/external/sgbuildex"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"time"

	"github.com/google/uuid"
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

// GetAuthorisations returns the currently stored authorisations for a user
func (s *PitstopService) GetAuthorisations(ctx context.Context, userID string) ([]*domain.PitstopAuthorisation, error) {
	return s.pitstopRepo.GetAuthorisationsByUser(ctx, userID)
}

// SyncConfig fetches the newest configs from the Pitstop API and upserts them
func (s *PitstopService) SyncConfig(ctx context.Context, userID string) error {
	// 1. Fetch from Pitstop API
	cfgResponse, err := s.pitstopClient.FetchConfig(ctx)
	if err != nil {
		return fmt.Errorf("pitstop API fetch failed: %w", err)
	}

	var authorisations []*domain.PitstopAuthorisation
	now := time.Now()

	// 2. Map JSON Response to Domain Entities
	for _, produce := range cfgResponse.Data.Produces {
		datasetID := produce.ID
		datasetName := produce.Name

		for _, to := range produce.To {
			regulatorID := to.ID
			regulatorName := to.Name

			for _, behalf := range to.OnBehalfOf {
				mainconID := behalf.ID
				mainconName := behalf.Name

				// Create the domain entity
				auth := &domain.PitstopAuthorisation{
					PitstopAuthID: "pitstop_auth_" + uuid.NewString()[:8], // fallback unique naming
					DatasetID:     datasetID,
					DatasetName:   datasetName,
					UserID:        userID,
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
