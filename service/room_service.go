// service/room_service.go
package service

import (
	"context"
	"log"

	"chat-websocket/redis"
)

// RoomService defines core room business logic.
type RoomService interface {
	BroadcastToRoom(ctx context.Context, roomName, message string) error
}

type roomServiceImpl struct {
	pubSubRepo redis.PubSubRepository
}

// NewRoomService creates a new RoomService instance.
func NewRoomService(pubSubRepo redis.PubSubRepository) RoomService {
	return &roomServiceImpl{
		pubSubRepo: pubSubRepo,
	}
}

func (r *roomServiceImpl) BroadcastToRoom(ctx context.Context, roomName, message string) error {
	err := r.pubSubRepo.Publish(ctx, roomName, roomName+"|"+message)
	if err != nil {
		log.Printf("Failed to broadcast to room %s: %v", roomName, err)
		return err
	}
	return nil
}
