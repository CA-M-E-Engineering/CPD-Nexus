package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

// ExternalSubmitter abstracts the external API client for submitting payloads.
// This allows the service layer to remain decoupled from the concrete sgbuildex adapter.
type ExternalSubmitter interface {
	// SubmitManpowerUtilization submits a batch of attendance records to the external API.
	SubmitManpowerUtilization(ctx context.Context, repo SubmissionRepository, settings *domain.SystemSettings, rows []domain.AttendanceRow) (int, error)

	// FetchPitstopConfig retrieves the latest pitstop authorisation configuration.
	FetchPitstopConfig(ctx context.Context) (*PitstopConfigResponse, error)
}

// PitstopConfigResponse is the parsed response from the Pitstop config endpoint.
type PitstopConfigResponse struct {
	Produces []PitstopProduceConfig
}

// PitstopProduceConfig represents a single dataset config returned by the Pitstop API.
type PitstopProduceConfig struct {
	ID   string
	Name string
	To   []PitstopRegulatorConfig
}

// PitstopRegulatorConfig represents a regulator (BCA/HDB/LTA) within a produce config.
type PitstopRegulatorConfig struct {
	ID         string
	Name       string
	OnBehalfOf []PitstopOnBehalfConfig
}

// PitstopOnBehalfConfig represents a contractor acting on behalf of the submitter.
type PitstopOnBehalfConfig struct {
	ID   string
	Name string
}
