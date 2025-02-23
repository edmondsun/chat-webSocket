// service/message_service.go
package service

import (
	"context"
	"log"

	"chat-websocket/redis"
)

// MessageService defines core message business logic.
type MessageService interface {
	SaveMessage(ctx context.Context, roomName, message string) error
	BroadcastMessage(ctx context.Context, roomName, message string) error
}

type messageServiceImpl struct {
	pubSubRepo redis.PubSubRepository
}

// NewMessageService creates a new MessageService instance.
func NewMessageService(pubSubRepo redis.PubSubRepository) MessageService {
	return &messageServiceImpl{
		pubSubRepo: pubSubRepo,
	}
}

func (m *messageServiceImpl) SaveMessage(ctx context.Context, roomName, message string) error {
	// Here you could extend logic to save the message into a database.
	log.Printf("[MessageService] Saving message for room=%s: %s", roomName, message)
	return nil
}

func (m *messageServiceImpl) BroadcastMessage(ctx context.Context, roomName, message string) error {
	err := m.pubSubRepo.Publish(ctx, roomName, roomName+"|"+message)
	if err != nil {
		log.Printf("Failed to broadcast message to room %s: %v", roomName, err)
		return err
	}
	return nil
}
