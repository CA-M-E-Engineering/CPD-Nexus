package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Meta contains common request/response metadata
type Meta struct {
	RequestID string `json:"request_id"`
	SentAt    string `json:"sent_at"`
}

// Message is the standard envelope for all bridge communication
type Message struct {
	Meta    Meta            `json:"meta"`
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// Handler interface for bridge requests
type Handler interface {
	Handle(ctx context.Context, msg Message) (*Message, error)
}

// NewRequest creates a standardized message envelope for outgoing requests
func NewRequest(action string, payload interface{}) (Message, error) {
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return Message{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return Message{
		Meta: Meta{
			RequestID: fmt.Sprintf("req-%s", time.Now().Format("20060102150405")),
			SentAt:    time.Now().Format(time.RFC3339),
		},
		Action:  action,
		Payload: rawPayload,
	}, nil
}
