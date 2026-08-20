package session

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/persistence/database/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/messages"
)

type SessionManager struct {
	sessionRepo *repositories.SessionRepository
	messageRepo *repositories.MessageRepository
}

func (m *SessionManager) CreateSession() (string, error) {
	ctx := context.Background()
	session, err := m.sessionRepo.Create(ctx)
	if err != nil {
		SessionLogger.Error("failed to create session", "error", err)
		return "", custom_errors.NewSessionError("failed to create session")
	}

	return session.ID.String(), nil
}

func (m *SessionManager) SaveSession(sessionID string, list messages.ListMessage) error {
	ctx := context.Background()
	mapper := &mappers.MessageMapper{}
	models, err := helpers.Map(list, func(t messages.Message) (*models.Message, error) {
		return mapper.ToModel(t, sessionID)
	})
	if err != nil {
		SessionLogger.With("SessionID", sessionID).Error("failed to map messages", "error", err)
		return custom_errors.NewSessionError("failed save session")
	}

	err = m.messageRepo.SaveMessages(ctx, models)
	if err != nil {
		SessionLogger.With("SessionID", sessionID).Error("failed save session", "error", err)
		return custom_errors.NewSessionError("failed save session")
	}

	return nil
}

func (m *SessionManager) GetMessagesBySessionID(id string) ([]messages.Message, error) {
	ctx := context.Background()
	mapper := &mappers.MessageMapper{}
	list, err := m.messageRepo.GetMessageBySessionID(ctx, id)
	if err != nil {
		SessionLogger.With("SessionID", id).Error("failed to get messages", "error", err)
		return nil, custom_errors.NewSessionError("failed to get messages")
	}

	return helpers.Map(list, func(t models.Message) (messages.Message, error) {
		return mapper.ToDomain(t)
	})
}

func (m *SessionManager) GetOrCreateSession(sessionID string) ([]messages.Message, error) {
	var targetSessionID string

	if sessionID == "" {
		newID, err := m.CreateSession()
		if err != nil {
			return nil, err
		}
		targetSessionID = newID
	} else {
		session, err := m.sessionRepo.FindByID(sessionID)
		if err != nil || session == nil {
			newID, err := m.CreateSession()
			if err != nil {
				return nil, err
			}
			targetSessionID = newID
		} else {
			targetSessionID = session.ID.String()
		}
	}

	return m.GetMessagesBySessionID(targetSessionID)
}

func (m *SessionManager) GetOrCreateSessionByUserID(userID string) (string, error) {
	ctx := context.Background()
	session, err := m.sessionRepo.FindByUserID(userID)
	if err == nil && session != nil {
		return session.ID.String(), err
	}

	newSession := &models.Session{
		UserID: userID,
	}

	createdSession, err := m.sessionRepo.CreateWithModel(ctx, newSession)
	if err != nil {
		SessionLogger.Error("failed to create session for user", "user_id", userID, "error", err)
		return "", custom_errors.NewSessionError("failed to create session")
	}

	return createdSession.ID.String(), nil
}

func (m *SessionManager) ResetSessionByID(sessionID string) error {
	ctx := context.Background()
	err := m.messageRepo.DeleteMessagesBySessionID(ctx, sessionID)
	if err != nil {
		SessionLogger.Error("failed to reset session", "session_id", sessionID, "error", err)
		return custom_errors.NewSessionError("failed to reset session")
	}

	return nil
}

// @Injectable
func NewSessionManager(sessionRepo *repositories.SessionRepository, messageRepo *repositories.MessageRepository) *SessionManager {
	return &SessionManager{
		sessionRepo: sessionRepo,
		messageRepo: messageRepo,
	}
}
