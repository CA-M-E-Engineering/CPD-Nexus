package sgbuildex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client represents the Ingress API client for SGBuildex
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	JWT        string // store JWT here after generation
}

// NewClient creates a new Ingress API client
func NewClient(baseURL string) *Client {
	client := &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
	client.JWT = client.GenerateJWT()
	return client
}

// GenerateJWT generates a JWT token for authentication
func (c *Client) GenerateJWT() string {
	// TODO: Implement real JWT creation using secret/key
	return "dummy-token"
}

// PostJSON sends a JSON payload to the specified endpoint
func (c *Client) PostJSON(endpoint string, payload any) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.JWT)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}
