package sgbuildex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client represents the Ingress API client for SGBuildex
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string // store API Key here
}

// NewClient creates a new Ingress API client
func NewClient(baseURL string) *Client {
	client := &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
	client.APIKey = client.FetchAPIKey()
	return client
}

// FetchAPIKey gets the SGTRADEX API key from environment
func (c *Client) FetchAPIKey() string {
	return strings.TrimSpace(os.Getenv("SGTRADEX_API_KEY"))
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
	if c.APIKey != "" {
		req.Header.Set("SGTRADEX-API-KEY", c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}
