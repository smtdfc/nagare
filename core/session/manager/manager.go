package manager

import (
	"context"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/shared/message"
)

type SessionManager struct {
	logger        *logger.BaseLogger
	sessionRepo   *repositories.SessionRepository
	messageRepo   *repositories.MessageRepository
	sessionMapper *mappers.SessionMapper
	messageMapper *mappers.MessageMapper
}

func (s *SessionManager) CreateUserSession(ctx context.Context, title string) (*session.SessionInfo, error) {
	sessionInfo := &session.SessionInfo{
		Title:     title,
		OwnerID:   uuid.Nil,
		OwnerType: session.USER,
		IsArchive: false,
	}

	newSession, err := s.sessionRepo.Create(ctx, s.sessionMapper.ToEntity(sessionInfo))
	if err != nil {
		return nil, custom_errors.ErrCreateUserSessionFailed
	}

	return s.sessionMapper.ToDomain(newSession), nil
}

func (s *SessionManager) GetListUserSession(ctx context.Context) ([]*session.SessionInfo, error) {
	sessions, err := s.sessionRepo.FindByOwnerType(ctx, session.USER.ToString())
	if err != nil {
		return nil, custom_errors.ErrGetUserSessionFailed
	}

	return s.sessionMapper.ToDomains(sessions), nil
}

func (s *SessionManager) GetUserSession(ctx context.Context, sessionID string) (*session.SessionInfo, error) {
	userSession, err := s.sessionRepo.FindUserSession(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to get user session", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetUserSessionFailed
	}

	if userSession == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	return s.sessionMapper.ToDomain(userSession), nil
}

func (s *SessionManager) GetUserChatHistory(ctx context.Context, sessionID string) (message.ListMessage, error) {
	chatSession, err := s.sessionRepo.FindUserSessionWithMessages(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to get user chat history", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetUserSessionFailed
	}

	if chatSession == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	domains, err := s.messageMapper.ToDomains(chatSession.Messages)
	if err != nil {
		s.logger.Error("failed to get user chat history", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetChatHistoryFailed
	}

	return domains, nil
}

func (s *SessionManager) SaveHistory(ctx context.Context, sessionID string, pendingMessage message.ListMessage) error {
	chatSession, err := s.sessionRepo.FindUserSessionWithMessages(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to save chat history", "session_id", sessionID, "err", err)
		return custom_errors.ErrSaveUserSessionFailed
	}

	if chatSession == nil {
		return custom_errors.ErrSessionNotFound
	}

	entities, err := s.messageMapper.ToEntities(pendingMessage, sessionID)
	if err != nil {
		s.logger.Error("failed to save chat history", "session_id", sessionID, "err", err)
		return custom_errors.ErrSaveUserSessionFailed
	}

	err = s.messageRepo.CreateBatch(ctx, entities, 200)
	if err != nil {
		s.logger.Error("failed to save chat history", "session_id", sessionID, "err", err)
		return custom_errors.ErrSaveUserSessionFailed
	}

	return nil
}

// @Injectable
func NewSessionManager(
	logger *logger.BaseLogger,
	sessionRepo *repositories.SessionRepository,
	sessionMapper *mappers.SessionMapper,
	messageRepo *repositories.MessageRepository,
	messageMapper *mappers.MessageMapper,
) *SessionManager {
	return &SessionManager{
		sessionRepo:   sessionRepo,
		sessionMapper: sessionMapper,
		messageRepo:   messageRepo,
		logger:        logger,
		messageMapper: messageMapper,
	}
}
