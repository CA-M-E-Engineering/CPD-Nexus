package bridge

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// Transport handles the low-level WebSocket connection and state
type Transport struct {
	url   string
	token string
	conn  *websocket.Conn
	mu    sync.Mutex
}

func NewTransport(bridgeURL, token string) *Transport {
	return &Transport{
		url:   bridgeURL,
		token: token,
	}
}

// Connect dial the bridge and maintains the connection
func (t *Transport) Connect() error {
	u, err := url.Parse(t.url)
	if err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	t.mu.Lock()
	t.conn = conn
	t.mu.Unlock()

	log.Printf("Transport: Connected to Bridge at %s", t.url)
	return nil
}

// Write sends a JSON message to the bridge
func (t *Transport) Write(msg Message) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return fmt.Errorf("transport not connected")
	}

	if msg.Meta.AuthToken == "" && t.token != "" {
		msg.Meta.AuthToken = t.token
	}

	return t.conn.WriteJSON(msg)
}

// Read waits for a single JSON message from the bridge
func (t *Transport) Read() (Message, error) {
	var msg Message

	if t.conn == nil {
		return msg, fmt.Errorf("transport not connected")
	}

	err := t.conn.ReadJSON(&msg)
	return msg, err
}

// Close explicitly closes the connection
func (t *Transport) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
}

// IsConnected returns true if the connection is active
func (t *Transport) IsConnected() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.conn != nil
}
