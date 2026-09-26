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

func (m *MessageRepository) CreateBatch(ctx context.Context, messages []*entities.Message, batchSize int) error {
	if len(messages) == 0 {
		return nil
	}

	if batchSize <= 0 {
		batchSize = 100
	}

	err := m.db.WithContext(ctx).CreateInBatches(messages, batchSize).Error
	if err != nil {
		m.logger.Error("Failed to create messages in batch", "count", len(messages), "error", err)
		return fmt.Errorf("failed to create messages in batch: %w", err)
	}

	return nil
}

func (m *MessageRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*entities.Message, error) {
	var messages []*entities.Message

	err := m.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC, rowid ASC").
		Find(&messages).Error

	if err != nil {
		m.logger.Error("Failed to get messages by session ID", "session", sessionID, "error", err)
		return nil, fmt.Errorf("failed to get messages by session ID: %w", err)
	}

	return messages, nil
}

func (m *MessageRepository) FindBySessionIDCursor(ctx context.Context, sessionID string, beforeID string, limit int) ([]*entities.Message, string, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	type messagePageRow struct {
		entities.Message
	}
	var rows []messagePageRow

	query := m.db.WithContext(ctx).
		Table("messages").
		Select("messages.*").
		Where("session_id = ?", sessionID)

	if beforeID != "" {
		query = query.Where(
			"(created_at, rowid) < (SELECT created_at, rowid FROM messages WHERE id = ? AND session_id = ?)",
			beforeID,
			sessionID,
		)
	}

	err := query.
		Order("created_at DESC, rowid DESC").
		Limit(limit + 1).
		Scan(&rows).Error
	if err != nil {
		m.logger.Error("Failed to get message page by session ID", "session", sessionID, "error", err)
		return nil, "", fmt.Errorf("failed to get message page by session ID: %w", err)
	}

	if len(rows) > limit {
		rows = rows[:limit]
	}

	messages := make([]*entities.Message, len(rows))
	for index := range rows {
		messages[index] = &rows[index].Message
	}

	for left, right := 0, len(messages)-1; left < right; left, right = left+1, right-1 {
		messages[left], messages[right] = messages[right], messages[left]
	}

	if len(messages) == 0 {
		return messages, "", nil
	}

	return messages, messages[0].ID.String(), nil
}

func (m *MessageRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	err := m.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Delete(&entities.Message{}).Error

	if err != nil {
		m.logger.Error("Failed to delete messages by session ID", "session", sessionID, "error", err)
		return fmt.Errorf("failed to delete messages by session ID: %w", err)
	}

	return nil
}

// @Injectable
func NewMessageRepository(db *gorm.DB, logger *logger.BaseLogger) *MessageRepository {
	return &MessageRepository{
		db:     db,
		logger: logger.With("module", "messages-repository"),
	}
}
