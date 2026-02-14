package sgbuildex

import (
	"context"
	"fmt"
)

// Health checks the API health endpoint
func (c *Client) Health(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v1/health", c.BaseURL)
	_, err := c.doRequest(ctx, "GET", url, nil, c.GenerateJWT())
	return err
}

// TODO: Add Config, UploadFiles, DownloadFile helpers
