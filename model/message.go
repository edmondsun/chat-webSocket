package model

import "time"

// Message represents a chat message structure.
type Message struct {
	ID        int64     `json:"id,omitempty"`
	SenderID  string    `json:"sender_id,omitempty"`
	RoomID    string    `json:"room_id,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`

	// For WebSocket actions:
	Action string `json:"action,omitempty"` // join, leave, message
}
