package sgbuildex

import (
	"context"
	"encoding/json"
	"fmt"
)

// PitstopConfigResponse defines the structure returned by /api/v1/config
type PitstopConfigResponse struct {
	Timestamp string `json:"timestamp"`
	Data      struct {
		Produces []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
			To   []struct {
				Name       string `json:"name"`
				ID         string `json:"id"`
				OnBehalfOf []struct {
					Name string `json:"name"`
					ID   string `json:"id"`
				} `json:"on_behalf_of"`
			} `json:"to"`
		} `json:"produces"`
	} `json:"data"`
}

// FetchConfig retrieves the routing configuration mapping from Pitstop API
func (c *Client) FetchConfig(ctx context.Context) (*PitstopConfigResponse, error) {
	url := fmt.Sprintf("%s/api/v1/config", c.PitstopURL)

	data, err := c.doRequest(ctx, "GET", url, nil, c.FetchAPIKey())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pitstop config: %w", err)
	}

	var response PitstopConfigResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse pitstop config json: %w", err)
	}

	return &response, nil
}
