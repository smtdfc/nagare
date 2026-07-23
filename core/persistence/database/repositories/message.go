package repositories

import (
	"context"

	"github.com/smtdfc/nagare/core/persistence/database"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func (r *MessageRepository) GetMessageBySessionID(ctx context.Context, id string) ([]models.Message, error) {
	var messages []models.Message
	messages, err := gorm.G[models.Message](r.db).
		Where("session_id = ?", id).
		Order("created_at ASC").
		Find(ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func NewMessageRepository() *MessageRepository {
	db, _ := database.GetDatabase()
	return &MessageRepository{
		db: db,
	}
}
