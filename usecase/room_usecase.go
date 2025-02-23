// usecase/room_usecase.go
package usecase

import (
	"context"
	"log"
	"strings"
	"sync"

	"chat-websocket/model"
	"chat-websocket/redis"
	"github.com/gorilla/websocket"
)

// RoomUseCase manages room operations such as join, leave, and local broadcasting.
type RoomUseCase struct {
	pubSubRepo redis.PubSubRepository
	rooms      map[string]*model.Room
	mutex      sync.RWMutex
}

// NewRoomUseCase creates a new RoomUseCase instance.
func NewRoomUseCase(pubSubRepo redis.PubSubRepository) *RoomUseCase {
	return &RoomUseCase{
		pubSubRepo: pubSubRepo,
		rooms:      make(map[string]*model.Room),
	}
}

// startPubSubListener listens for cross-server messages via Redis Pub/Sub.
func (uc *RoomUseCase) startPubSubListener(roomName string) {
	go uc.pubSubRepo.Subscribe(context.Background(), roomName, func(payload []byte) {
		parts := strings.SplitN(string(payload), "|", 2)
		if len(parts) != 2 {
			log.Printf("[RoomUseCase] Invalid broadcast for room %s: %s", roomName, payload)
			return
		}
		uc.BroadcastToLocalRoom(roomName, parts[1])
	})
	log.Printf("[RoomUseCase] Started PubSub listener for room %s", roomName)
}

// JoinRoom adds a client to a room and publishes a join message.
func (uc *RoomUseCase) JoinRoom(ctx context.Context, client *model.Client, roomName string) {
	uc.mutex.Lock()
	room, exists := uc.rooms[roomName]
	if !exists {
		room = &model.Room{
			Name:    roomName,
			Clients: make(map[string]*model.ClientConn),
		}
		uc.rooms[roomName] = room
		uc.startPubSubListener(roomName)
	}
	room.Mutex.Lock()
	room.Clients[client.ID] = &model.ClientConn{
		ID:   client.ID,
		Conn: client,
	}
	room.Mutex.Unlock()
	uc.mutex.Unlock()

	log.Printf("[RoomUseCase] Client %s joined room %s", client.ID, roomName)
	_ = uc.pubSubRepo.Publish(ctx, roomName, roomName+"|"+client.ID+" joined the room")
}

// LeaveRoom removes a client from the specified room and publishes a leave message.
func (uc *RoomUseCase) LeaveRoom(ctx context.Context, clientID, roomName string) {
	uc.mutex.RLock()
	room, exists := uc.rooms[roomName]
	uc.mutex.RUnlock()
	if !exists {
		log.Printf("[RoomUseCase] Room %s does not exist", roomName)
		return
	}

	room.Mutex.Lock()
	delete(room.Clients, clientID)
	empty := len(room.Clients) == 0
	room.Mutex.Unlock()

	if empty {
		uc.mutex.Lock()
		delete(uc.rooms, roomName)
		uc.mutex.Unlock()
		uc.pubSubRepo.Unsubscribe(ctx, roomName)
	}

	log.Printf("[RoomUseCase] Client %s left room %s", clientID, roomName)
	_ = uc.pubSubRepo.Publish(ctx, roomName, roomName+"|"+clientID+" left the room")
}

// RemoveClient removes a client from all rooms.
func (uc *RoomUseCase) RemoveClient(ctx context.Context, clientID string) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	for roomName, room := range uc.rooms {
		room.Mutex.Lock()
		if _, exists := room.Clients[clientID]; exists {
			delete(room.Clients, clientID)
			log.Printf("[RoomUseCase] Client %s removed from room %s", clientID, roomName)
			if len(room.Clients) == 0 {
				delete(uc.rooms, roomName)
				uc.pubSubRepo.Unsubscribe(ctx, roomName)
			}
			_ = uc.pubSubRepo.Publish(ctx, roomName, roomName+"|"+clientID+" left the room")
		}
		room.Mutex.Unlock()
	}
}

// BroadcastMessage broadcasts a message to all servers via Redis.
func (uc *RoomUseCase) BroadcastMessage(ctx context.Context, roomName, message string) {
	if err := uc.pubSubRepo.Publish(ctx, roomName, roomName+"|"+message); err != nil {
		log.Printf("[RoomUseCase] Failed to broadcast message to room %s: %v", roomName, err)
	}
}

// BroadcastToLocalRoom sends a message to all clients in the room on the local server.
func (uc *RoomUseCase) BroadcastToLocalRoom(roomName, message string) {
	uc.mutex.RLock()
	room, exists := uc.rooms[roomName]
	uc.mutex.RUnlock()
	if !exists {
		log.Printf("[RoomUseCase] Room %s does not exist for local broadcast.", roomName)
		return
	}

	room.Mutex.RLock()
	defer room.Mutex.RUnlock()

	for _, conn := range room.Clients {
		go func(cc *model.ClientConn) {
			cc.Conn.Mutex.Lock()
			defer cc.Conn.Mutex.Unlock()
			if cc.Conn.Conn == nil {
				return
			}
			if err := cc.Conn.Conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Printf("Failed to send message to client %s: %v", cc.ID, err)
			}
		}(conn)
	}
}
