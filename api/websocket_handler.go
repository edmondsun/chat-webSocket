// api/websocket_handler.go
package api

import (
	"chat-websocket/model"
	"chat-websocket/pkg/metrics"
	"chat-websocket/usecase"
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections and incoming messages.
type WebSocketHandler struct {
	RoomUseCase    *usecase.RoomUseCase
	MessageUseCase *usecase.MessageUseCase
	Upgrader       websocket.Upgrader
}

// NewWebSocketHandler creates a new WebSocketHandler instance.
func NewWebSocketHandler(roomUseCase *usecase.RoomUseCase, messageUseCase *usecase.MessageUseCase) *WebSocketHandler {
	return &WebSocketHandler{
		RoomUseCase:    roomUseCase,
		MessageUseCase: messageUseCase,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development; adjust in production.
				return true
			},
		},
	}
}

// HandleConnection upgrades the HTTP connection to a WebSocket and processes messages.
func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}
	defer func() {
		log.Println("Closing WebSocket connection.")
		conn.Close()
	}()

	// Set read deadline for heartbeat.
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	senderID := r.URL.Query().Get("sender_id")
	if senderID == "" {
		log.Println("WebSocket connection rejected: sender_id is missing.")
		conn.Close()
		return
	}
	log.Printf("WebSocket connection request from: %s, sender_id: %s", conn.RemoteAddr(), senderID)

	client := &model.Client{
		ID:       conn.RemoteAddr().String(),
		Conn:     conn,
		SenderID: senderID,
	}

	defer func() {
		h.RoomUseCase.RemoveClient(context.Background(), client.ID)
		log.Printf("Client disconnected: %s\n", client.ID)
	}()

	// Channel to receive messages from the connection non-blockingly.
	messageChan := make(chan []byte, 50)
	go h.readMessages(conn, messageChan)

	for msg := range messageChan {
		var incoming model.Message
		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Printf("Invalid message format: %v\n", err)
			continue
		}
		h.handleMessage(client, incoming)
	}
}

// readMessages reads messages from the WebSocket connection asynchronously.
// When encountering errors (e.g. i/o timeout), it increases a counter and exits after reaching a threshold.
func (h *WebSocketHandler) readMessages(conn *websocket.Conn, messageChan chan<- []byte) {
	defer close(messageChan)
	timeoutCount := 0
	const maxTimeouts = 5 // Exit reading after 5 timeouts

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			metrics.ReadErrors.Inc() // Increment error counter
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				timeoutCount++
				log.Printf("WebSocket read timeout (%d/%d): %v", timeoutCount, maxTimeouts, err)
				if timeoutCount >= maxTimeouts {
					log.Printf("Maximum timeout reached. Closing connection.")
					break
				}
				time.Sleep(1 * time.Minute)
				continue
			} else {
				log.Printf("WebSocket read error (non-timeout): %v", err)
				break
			}
		}
		// Reset timeout counter on successful read.
		timeoutCount = 0
		metrics.MessagesRead.Inc() // Increment messages read counter
		messageChan <- data
	}
}

// handleMessage processes the incoming message based on its action.
// It checks if the RoomID is provided; if empty, it logs an error and ignores the message.
func (h *WebSocketHandler) handleMessage(client *model.Client, msg model.Message) {
	// Validate that RoomID is not empty.
	if msg.RoomID == "" {
		log.Printf("Error: RoomID is empty in message from client %s", client.ID)
		return
	}

	msg.SenderID = client.SenderID

	switch msg.Action {
	case "join":
		h.RoomUseCase.JoinRoom(context.Background(), client, msg.RoomID)
	case "leave":
		h.RoomUseCase.LeaveRoom(context.Background(), client.ID, msg.RoomID)
	case "message":
		// Process the message: save to DB and broadcast.
		h.MessageUseCase.ProcessMessage(context.Background(), msg)
	default:
		log.Printf("Unknown action: %s", msg.Action)
	}
}
