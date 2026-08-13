package session

import (
	"context"

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
		return "", err
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
		return err
	}

	err = m.messageRepo.SaveMessages(ctx, models)
	if err != nil {
		return err
	}

	return nil
}

func (m *SessionManager) GetMessagesBySessionID(id string) ([]messages.Message, error) {
	ctx := context.Background()
	mapper := &mappers.MessageMapper{}
	list, err := m.messageRepo.GetMessageBySessionID(ctx, id)
	if err != nil {
		return nil, err
	}

	return helpers.Map(list, func(t models.Message) (messages.Message, error) {
		return mapper.ToDomain(t)
	})
}

func NewSessionManager(sessionRepo *repositories.SessionRepository, messageRepo *repositories.MessageRepository) *SessionManager {
	return &SessionManager{
		sessionRepo: sessionRepo,
		messageRepo: messageRepo,
	}
}
