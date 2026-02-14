package sgbuildex

import (
	"context"
	"fmt"
)

// Push sends data to a Push API
func (c *Client) PushEvent(ctx context.Context, dataElementID string, payload any, participants []ParticipantWrapper, onBehalfOf []OnBehalfWrapper) error {
	req := PushRequest{
		Participants: participants,
		Payload:      []any{payload},
		OnBehalfOf:   onBehalfOf,
	}

	url := fmt.Sprintf("%s/api/v1/data/push/%s", c.BaseURL, dataElementID)
	_, err := c.doRequest(ctx, "POST", url, req, c.GenerateJWT())
	return err
}

// TODO: Add Store, Pull, Forward functions
