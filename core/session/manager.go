package session

import (
	"context"

	"github.com/smtdfc/nagare/core/persistence/database/repositories"
)

type SessionManager struct {
	repo *repositories.SessionRepository
}

func (m *SessionManager) CreateSession() (string, error) {
	ctx := context.Background()
	session, err := m.repo.Create(ctx)
	if err != nil {
		return "", err
	}

	return session.ID.String(), nil
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		repo: repositories.NewSessionRepository(),
	}
}
