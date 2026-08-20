package repositories

import (
	"context"

	"github.com/smtdfc/nagare/core/persistence"
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
		persistence.PersistenceLogger.Error("Failed to get messages by session ID", "error", err)
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) DeleteMessagesBySessionID(ctx context.Context, id string) error {
	_, err := gorm.G[models.Message](r.db).
		Where("session_id = ?", id).Delete(ctx)
	if err != nil {
		persistence.PersistenceLogger.Error("Failed to delete messages by session ID", "error", err)
		return err
	}

	return nil
}

func (r *MessageRepository) SaveMessages(ctx context.Context, messages []*models.Message) error {
	if len(messages) == 0 {
		return nil
	}

	batchSize := 500

	err := r.db.WithContext(ctx).CreateInBatches(messages, batchSize).Error
	if err != nil {
		persistence.PersistenceLogger.Error("Failed to save messages", "error", err)
		return err
	}

	return nil
}

// @Injectable
func NewMessageRepository() *MessageRepository {
	db, _ := database.GetDatabase()
	return &MessageRepository{
		db: db,
	}
}
