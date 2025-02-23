// repository/message_repository.go
package repository

import (
	"chat-websocket/model"

	"gorm.io/gorm"
)

// MessageRepository defines methods for accessing message data.
type MessageRepository interface {
	CreateMessage(msg *model.Message) error
	GetMessagesByRoom(room string) ([]model.Message, error)
}

// MysqlMessageRepository is the MySQL implementation of MessageRepository.
type MysqlMessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new instance of MysqlMessageRepository.
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &MysqlMessageRepository{db: db}
}

func (r *MysqlMessageRepository) CreateMessage(msg *model.Message) error {
	return r.db.Create(msg).Error
}

func (r *MysqlMessageRepository) GetMessagesByRoom(room string) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Where("room = ?", room).Order("created_at ASC").Find(&messages).Error
	return messages, err
}
