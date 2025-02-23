// repository/client_repository.go
package repository

import (
	"chat-websocket/model"
	"fmt"

	"gorm.io/gorm"
)

// ClientRepository defines methods for client data access.
type ClientRepository interface {
	CreateClient(client *model.Client) error
	GetClientByID(clientID string) (*model.Client, error)
	UpdateClient(client *model.Client) error
}

// MysqlClientRepository is the MySQL implementation of ClientRepository.
type MysqlClientRepository struct {
	db *gorm.DB
}

// NewClientRepository creates a new instance of MysqlClientRepository.
func NewClientRepository(db *gorm.DB) ClientRepository {
	return &MysqlClientRepository{db: db}
}

func (r *MysqlClientRepository) CreateClient(client *model.Client) error {
	return r.db.Create(client).Error
}

func (r *MysqlClientRepository) GetClientByID(clientID string) (*model.Client, error) {
	var client model.Client
	if err := r.db.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return nil, fmt.Errorf("client not found: %s", clientID)
	}
	return &client, nil
}

func (r *MysqlClientRepository) UpdateClient(client *model.Client) error {
	return r.db.Save(client).Error
}
