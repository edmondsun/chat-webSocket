// model/client.go
package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected user's WebSocket session along with optional user data.
type Client struct {
	ID    string          // Unique identifier (e.g., remote address)
	Conn  *websocket.Conn // WebSocket connection
	Mutex sync.Mutex      // To protect write operations

	// Optional database fields:
	ClientID string
	Email    string
	Name     string
	Status   int
	SenderID string
}
