package repositories

import (
	"context"
	"fmt"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db     *gorm.DB
	logger *logger.BaseLogger
}

func (m *MessageRepository) GetListMessageBySessionID(ctx context.Context, sessionID string) ([]*entities.Message, error) {
	var messages []*entities.Message

	err := m.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error

	if err != nil {
		m.logger.Error("Failed to get messages by session ID", "session", sessionID, "error", err)
		return nil, fmt.Errorf("failed to get messages by session ID: %w", err)
	}

	return messages, nil
}

// @Injectable
func NewMessageRepository(db *gorm.DB, logger *logger.BaseLogger) *MessageRepository {
	return &MessageRepository{
		db:     db,
		logger: logger.With("module", "message-repository"),
	}
}
