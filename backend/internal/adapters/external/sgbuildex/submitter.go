package sgbuildex

import (
	"context"
	"encoding/json"
	"fmt"
	"sgbuildex/internal/adapters/external/sgbuildex/payloads"
	"sgbuildex/internal/core/ports"
)

// Submittable defines the behavior for any payload that can be pushed to SGBuildex
type Submittable interface {
	DataElementID() string
	ToPushRequest(ctx context.Context) (*PushRequest, error)
	GetInternalID() string
}

// SubmitPayloads submits any submittable payloads to SGBuildex
// and updates the database with status, response payload, and errors via the repository.
func SubmitPayloads[T Submittable](ctx context.Context, repo ports.SubmissionRepository, client *Client, submittables []T) error {
	for _, s := range submittables {
		// Prepare Push Request
		pushReq, err := s.ToPushRequest(ctx)
		if err != nil {
			fmt.Printf("Failed to prepare push request for %s: %v\n", s.GetInternalID(), err)
			continue
		}

		fullJSON, err := json.MarshalIndent(pushReq, "", "  ")
		if err != nil {
			fmt.Printf("Failed to marshal full push request: %s\n", err)
		} else {
			fmt.Println("====================================================")
			fmt.Printf("SUBMISSION JSON PAYLOAD (%s):\n", s.DataElementID())
			fmt.Println(string(fullJSON))
			fmt.Println("====================================================")
		}

		// Push event
		err = client.PushEvent(ctx, s.DataElementID(), pushReq.Payload[0], pushReq.Participants, pushReq.OnBehalfOf)

		// Determine status
		status := "submitted"
		errorMessage := ""
		if err != nil {
			status = "failed"
			errorMessage = err.Error()
			fmt.Printf("Failed to submit payload: %s\n", err)
		} else {
			fmt.Println("Payload submitted successfully")
		}

		// 1. Log to central submission_logs table via Repo
		logErr := repo.LogSubmission(ctx, s.DataElementID(), s.GetInternalID(), status, string(fullJSON), errorMessage)
		if logErr != nil {
			fmt.Printf("Failed to write to submission_logs: %v\n", logErr)
		}

		// 2. Update specific source table if needed
		if s.DataElementID() == "manpower_utilization" {
			dbErr := repo.UpdateAttendanceStatus(ctx, s.GetInternalID(), status, string(fullJSON), errorMessage)
			if dbErr != nil {
				fmt.Printf("Failed to update status for %s: %s\n", s.GetInternalID(), dbErr)
			}
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
