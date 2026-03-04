package sgbuildex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client is the HTTP client for communicating with the SGBuildex / Pitstop API.
type Client struct {
	BaseURL    string
	PitstopURL string
	HTTPClient *http.Client
	APIKey     string
}

// NewClient creates a new Client and loads the API key from the environment.
func NewClient(baseURL, pitstopURL string) *Client {
	apiKey := strings.TrimSpace(os.Getenv("SGTRADEX_API_KEY"))
	if apiKey == "" {
		log.Printf("[SGBuildex] WARNING: SGTRADEX_API_KEY is not set — requests may be rejected.")
	}
	return &Client{
		BaseURL:    baseURL,
		PitstopURL: pitstopURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		APIKey:     apiKey,
	}
}

// PostJSON marshals payload as JSON and POSTs it to the given endpoint on the Pitstop server.
// The SGTRADEX-API-KEY header is set automatically if an API key is configured.
func (c *Client) PostJSON(endpoint string, payload any) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.PitstopURL, endpoint)

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.APIKey != "" {
		req.Header.Set("SGTRADEX-API-KEY", c.APIKey)
		req.Header.Set("x-api-key", c.APIKey)
	}

	log.Printf("[SGBuildex] POST %s", url)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}
