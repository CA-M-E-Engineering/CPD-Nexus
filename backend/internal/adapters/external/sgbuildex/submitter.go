package sgbuildex

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sgbuildex/internal/adapters/external/sgbuildex/payloads"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
	"time"
)

// Submittable defines the behavior for any payload that can be pushed to SGBuildex
type Submittable interface {
	DataElementID() string
	ToPushRequest(ctx context.Context) (*PushRequest, error)
	GetInternalID() string
}

// SubmitPayloads submissions any submittable payloads to SGBuildex in batches.
// It respects the MaxWorkersPerRequest and MaxPayloadSizeKB settings.
func SubmitPayloads[T Submittable](ctx context.Context, repo ports.SubmissionRepository, client *Client, settings *domain.SystemSettings, submittables []T) error {
	if len(submittables) == 0 {
		return nil
	}

	dataElementID := submittables[0].DataElementID()
	maxBatchSize := settings.MaxWorkersPerRequest
	if maxBatchSize <= 0 {
		maxBatchSize = 100
	}
	limitBytes := int(settings.MaxPayloadSizeKB) * 1024
	if limitBytes <= 0 {
		limitBytes = 256 * 1024
	}

	totalItems := len(submittables)
	for i := 0; i < totalItems; {
		var batchParticipants []ParticipantWrapper
		var batchPayload []any
		var batchOnBehalf []OnBehalfWrapper
		var batchIDs []string

		// Build the largest possible batch within limits
		for i < totalItems && len(batchIDs) < maxBatchSize {
			s := submittables[i]
			req, err := s.ToPushRequest(ctx)
			if err != nil {
				log.Printf("[SGBuildex] Failed to prepare %s for %s: %v", dataElementID, s.GetInternalID(), err)
				i++
				continue
			}

			// Preview if we add this item
			nextParticipants := append(batchParticipants, req.Participants...)
			nextPayload := append(batchPayload, req.Payload...)

			// Deduplicate OnBehalfOf
			seenUEN := make(map[string]bool)
			for _, ob := range batchOnBehalf {
				seenUEN[ob.ID] = true
			}
			nextOnBehalf := append([]OnBehalfWrapper{}, batchOnBehalf...)
			for _, ob := range req.OnBehalfOf {
				if !seenUEN[ob.ID] {
					nextOnBehalf = append(nextOnBehalf, ob)
					seenUEN[ob.ID] = true
				}
			}

			// Size check
			pushReq := &PushRequest{
				Participants: nextParticipants,
				Payload:      nextPayload,
				OnBehalfOf:   nextOnBehalf,
			}
			jsonBytes, _ := json.Marshal(pushReq)

			if len(jsonBytes) > limitBytes {
				if len(batchIDs) == 0 {
					// Single item above limit - skip and log
					log.Printf("[SGBuildex] CRITICAL: Single item for %s is already above size limit (%d > %d bytes). Skipping.", s.GetInternalID(), len(jsonBytes), limitBytes)
					i++
					continue
				}
				// Batch reached size limit, stop here and send what we have
				break
			}

			// Accept the item
			batchParticipants = nextParticipants
			batchPayload = nextPayload
			batchOnBehalf = nextOnBehalf
			batchIDs = append(batchIDs, s.GetInternalID())
			i++
		}

		if len(batchIDs) == 0 {
			continue
		}

		// Prepare final batch request
		finalReq := &PushRequest{
			Participants: batchParticipants,
			Payload:      batchPayload,
			OnBehalfOf:   batchOnBehalf,
		}
		fullJSON, _ := json.MarshalIndent(finalReq, "", "  ")

		log.Printf("[SGBuildex] Submitting batch of %d items for %s (Size: %d bytes)", len(batchIDs), dataElementID, len(fullJSON))

		// Execute submission
		resp, err := client.PostJSON(fmt.Sprintf("api/v1/data/push/%s", dataElementID), finalReq)

		status := "submitted"
		errorMessage := ""
		if err != nil {
			status = "failed"
			errorMessage = err.Error()
			log.Printf("[SGBuildex] Batch submission failed: %v", err)
		} else {
			resp.Body.Close()
			if resp.StatusCode >= 400 {
				status = "failed"
				errorMessage = fmt.Sprintf("HTTP %d", resp.StatusCode)
				log.Printf("[SGBuildex] Batch submission returned error: %s", errorMessage)
			}
		}

		// Update database for each individual item in the batch
		for _, id := range batchIDs {
			// Log to central logs
			repo.LogSubmission(ctx, dataElementID, id, status, string(fullJSON), errorMessage)

			// Update specific source table if needed
			if dataElementID == "manpower_utilization" {
				repo.UpdateAttendanceStatus(ctx, id, status, string(fullJSON), errorMessage)
			}
		}

		// Rate limiting safety: if we have more batches, wait a bit
		if i < totalItems && settings.MaxRequestsPerMinute > 0 {
			sleepDuration := time.Minute / time.Duration(settings.MaxRequestsPerMinute)
			time.Sleep(sleepDuration)
		}
	}

	return nil
}

// ManpowerUtilizationWrapper wraps the payload to implement Submittable
type ManpowerUtilizationWrapper struct {
	payloads.ManpowerUtilization
}

func (w ManpowerUtilizationWrapper) DataElementID() string {
	return "manpower_utilization"
}

func (w ManpowerUtilizationWrapper) GetInternalID() string {
	return w.InternalAttendanceID
}

func (w ManpowerUtilizationWrapper) ToPushRequest(ctx context.Context) (*PushRequest, error) {
	// Prepare Project Reference
	projectRef := ""
	if w.ProjectReferenceNumber != nil {
		projectRef = *w.ProjectReferenceNumber
	}

	// Use the pre-fetched PIC from the payload
	var participantOnBehalf *OnBehalfWrapper
	if w.InternalPICFIN != "" {
		participantOnBehalf = &OnBehalfWrapper{
			ID:   w.InternalPICFIN,
			Name: w.InternalPICName,
		}
	} else {
		fmt.Printf("Warning: No PIC found for project associated with worker %s\n", w.InternalWorkerID)
	}

	// Participants
	participants := []ParticipantWrapper{
		{
			ID:   w.PersonIDNo,
			Name: w.PersonName,
			Meta: &ParticipantMeta{
				DataRefID: projectRef,
			},
			OnBehalfOf: participantOnBehalf,
		},
	}

	// Request level OnBehalfOf (Participant's company and Main Contractor)
	onBehalfOf := []OnBehalfWrapper{
		{ID: w.PersonEmployerCompanyUEN},
	}
	if w.MainContractorCompanyUEN != nil && *w.MainContractorCompanyUEN != "" {
		onBehalfOf = append(onBehalfOf, OnBehalfWrapper{ID: *w.MainContractorCompanyUEN})
	}

	return &PushRequest{
		Participants: participants,
		Payload:      []any{w.ManpowerUtilization},
		OnBehalfOf:   onBehalfOf,
	}, nil
}

// ManpowerDistributionWrapper wraps the payload to implement Submittable
type ManpowerDistributionWrapper struct {
	payloads.ManpowerDistribution
}

func (w ManpowerDistributionWrapper) DataElementID() string {
	return "manpower_distribution"
}

func (w ManpowerDistributionWrapper) GetInternalID() string {
	return fmt.Sprintf("%s_%s", w.SubmissionMonth, w.OffsiteFabricatorCompanyUEN)
}

func (w ManpowerDistributionWrapper) ToPushRequest(ctx context.Context) (*PushRequest, error) {
	// Request level OnBehalfOf for the fabricator
	onBehalfOf := []OnBehalfWrapper{
		{ID: w.OffsiteFabricatorCompanyUEN},
	}

	return &PushRequest{
		Participants: []ParticipantWrapper{}, // Initialize to empty slice for [] in JSON
		Payload:      []any{w.ManpowerDistribution},
		OnBehalfOf:   onBehalfOf,
	}, nil
}
