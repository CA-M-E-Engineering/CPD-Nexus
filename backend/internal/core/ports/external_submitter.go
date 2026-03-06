package ports

import (
	"context"
	"sgbuildex/internal/core/domain"
)

// PitstopOnBehalfConfig maps the onBehalfOf fields
type PitstopOnBehalfConfig struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PitstopRegulatorConfig maps the to fields
type PitstopRegulatorConfig struct {
	ID         string                  `json:"id"`
	Name       string                  `json:"name"`
	OnBehalfOf []PitstopOnBehalfConfig `json:"onBehalfOf"`
}

// PitstopProduceConfig maps the produces fields
type PitstopProduceConfig struct {
	ID   string                   `json:"id"`
	Name string                   `json:"name"`
	To   []PitstopRegulatorConfig `json:"to"`
}

// PitstopConfigResponse models the expected output layout from pitstop fetch APIs
type PitstopConfigResponse struct {
	Produces []PitstopProduceConfig `json:"produces"`
}

// ExternalSubmitter defines the interface for external Pitstop/SGBuildex submissions
type ExternalSubmitter interface {
	FetchPitstopConfig(ctx context.Context) (*PitstopConfigResponse, error)
	SubmitManpowerUtilization(ctx context.Context, repo SubmissionRepository, settings *domain.SystemSettings, rows []domain.AttendanceRow) (int, int, error)
}
