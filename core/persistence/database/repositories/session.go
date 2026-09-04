package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db     *gorm.DB
	logger *logger.BaseLogger
}

func (s *SessionRepository) GetSessionByID(ctx context.Context, id string) (*entities.Session, error) {
	var session entities.Session
	err := s.db.WithContext(ctx).
		Where("id = ? AND owner_type = ?", id).
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}

		s.logger.Error("failed to get session by ID", "id", id, "error", err)
		return nil, err
	}

	return &session, nil
}

func (s *SessionRepository) GetListSessionByOwnerType(ctx context.Context, ownerType string) ([]*entities.Session, error) {
	var sessions []*entities.Session

	err := s.db.WithContext(ctx).
		Where("owner_type = ?", ownerType).
		Find(&sessions).Error

	if err != nil {
		s.logger.Error("failed to get list session by owner type", "owner_type", ownerType, "error", err)
		return nil, err
	}

	return sessions, nil
}

func (s *SessionRepository) GetListSessionByOwnerID(ctx context.Context, ownerType string, ownerID string) ([]*entities.Session, error) {
	var sessions []*entities.Session

	err := s.db.WithContext(ctx).
		Where("owner_type = ? AND owner_id = ?", ownerType, ownerID).
		Find(&sessions).Error

	if err != nil {
		s.logger.Error("failed to get list session by owner id", "owner_id", ownerID, "owner_type", ownerType, "error", err)
		return nil, err
	}

	return sessions, nil
}

func (s *SessionRepository) GetListSessionByChannelID(ctx context.Context, ownerType string, ownerID string, channelID string) ([]*entities.Session, error) {
	var sessions []*entities.Session
	err := s.db.WithContext(ctx).
		Where("owner_type = ? AND owner_id = ? AND channel_id = ?", ownerType, ownerID, channelID).
		Find(&sessions).Error

	if err != nil {
		s.logger.Error("failed to get list session by owner id", "channel", channelID, "owner_id", ownerID, "owner_type", ownerType, "error", err)
		return nil, err
	}

	return sessions, nil
}

func (s *SessionRepository) CreateSession(ctx context.Context, session *entities.Session) (*entities.Session, error) {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}

	err := s.db.WithContext(ctx).Create(session).Error
	if err != nil {
		s.logger.Error("failed to create session", "error", err)
		return nil, err
	}

	return session, nil
}

func (s *SessionRepository) UpdateSession(ctx context.Context, session *entities.Session) error {
	err := s.db.WithContext(ctx).Save(session).Error
	if err != nil {
		s.logger.Error("failed to update session", "id", session.ID, "error", err)
		return err
	}

	return nil
}

func (s *SessionRepository) DeleteSession(ctx context.Context, id string) error {
	result := s.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Session{})
	if result.Error != nil {
		s.logger.Error("failed to delete session", "id", id, "error", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("session not found to delete")
	}

	return nil
}

func (s *SessionRepository) GetUserSession(ctx context.Context, sessionId string) (*entities.Session, error) {
	var session entities.Session
	err := s.db.WithContext(ctx).
		Where("owner_type = ? AND id = ?", "user", sessionId).
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		s.logger.Error("failed to get session", "error", err)
		return nil, err
	}

	return &session, nil
}

func (s *SessionRepository) GetUserSessionWithMessages(ctx context.Context, sessionId string) (*entities.Session, error) {
	var session entities.Session

	err := s.db.WithContext(ctx).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&session, "owner_type = ? AND id = ?", "user", sessionId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		s.logger.Error("failed to get session", "error", err)
		return nil, err
	}
	return &session, nil
}

// @Injectable
func NewSessionRepository(db *gorm.DB, logger *logger.BaseLogger) *SessionRepository {
	return &SessionRepository{
		db:     db,
		logger: logger.With("module", "session-repository"),
	}
}
