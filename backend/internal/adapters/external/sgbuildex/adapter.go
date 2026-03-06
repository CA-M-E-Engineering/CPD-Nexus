package sgbuildex

import (
	"context"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/ports"
)

// Ensure *Client implements ports.ExternalSubmitter at compile time.
var _ ports.ExternalSubmitter = (*Client)(nil)

// SubmitManpowerUtilization implements ports.ExternalSubmitter.
// It maps the domain AttendanceRows to ManpowerUtilization payloads and submits them.
func (c *Client) SubmitManpowerUtilization(ctx context.Context, repo ports.SubmissionRepository, settings *domain.SystemSettings, rows []domain.AttendanceRow) (int, int, error) {
	muResult := MapAttendanceToManpower(rows)

	failedCount := 0
	for id, errMsg := range muResult.Failures {
		repo.UpdateAttendanceStatus(ctx, id, "failed", "", errMsg)
		repo.LogSubmission(ctx, "manpower_utilization", id, "failed", "", errMsg)
		failedCount++
	}

	wrappers := make([]ManpowerUtilizationWrapper, len(muResult.Payloads))
	for i, p := range muResult.Payloads {
		wrappers[i] = ManpowerUtilizationWrapper{ManpowerUtilization: p}
	}
	submittedCount, err := SubmitPayloads(ctx, repo, c, settings, wrappers)
	return submittedCount, failedCount, err
}

// FetchPitstopConfig implements ports.ExternalSubmitter — wraps the concrete FetchConfig method
// and converts the response to the ports-level type (no concrete adapter types escape the boundary).
func (c *Client) FetchPitstopConfig(ctx context.Context) (*ports.PitstopConfigResponse, error) {
	resp, err := c.FetchConfig(ctx)
	if err != nil {
		return nil, err
	}

	portResp := &ports.PitstopConfigResponse{}
	for _, prod := range resp.Data.Produces {
		pp := ports.PitstopProduceConfig{ID: prod.ID, Name: prod.Name}
		for _, to := range prod.To {
			pr := ports.PitstopRegulatorConfig{ID: to.ID, Name: to.Name}
			for _, ob := range to.OnBehalfOf {
				pr.OnBehalfOf = append(pr.OnBehalfOf, ports.PitstopOnBehalfConfig{ID: ob.ID, Name: ob.Name})
			}
			pp.To = append(pp.To, pr)
		}
		portResp.Produces = append(portResp.Produces, pp)
	}
	return portResp, nil
}
