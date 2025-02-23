// usecase/message_usecase.go
package usecase

import (
	"chat-websocket/model"
	"chat-websocket/repository"
	"chat-websocket/service"
	"context"
	"log"
)

// MessageUseCase encapsulates higher-level message processing logic.
type MessageUseCase struct {
	MessageRepo    repository.MessageRepository
	MessageService service.MessageService
}

// NewMessageUseCase creates a new instance of MessageUseCase.
func NewMessageUseCase(repo repository.MessageRepository, service service.MessageService) *MessageUseCase {
	return &MessageUseCase{
		MessageRepo:    repo,
		MessageService: service,
	}
}

// ProcessMessage processes an incoming message: it saves the message to the DB and broadcasts it.
func (mu *MessageUseCase) ProcessMessage(ctx context.Context, msg model.Message) {
	// Save the message to the database.
	if err := mu.MessageRepo.CreateMessage(&msg); err != nil {
		log.Printf("[MessageUseCase] Failed to save message: %v\n", err)
	} else {
		log.Printf("[MessageUseCase] Message saved successfully.")
	}

	// Broadcast the message using MessageService.
	// Use msg.RoomID instead of msg.Room.
	if err := mu.MessageService.BroadcastMessage(ctx, msg.RoomID, msg.Content); err != nil {
		log.Printf("[MessageUseCase] Failed to broadcast message: %s: %v\n", msg.SenderID, err)
	} else {
		log.Printf("[MessageUseCase] Message broadcasted successfully: %s.", msg.SenderID)
	}
}
