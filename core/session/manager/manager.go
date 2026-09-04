package manager

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/shared/message"
)

const DefaultUserID = "00000000-0000-0000-0000-000000000000"

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
		OwnerID:   DefaultUserID,
		OwnerType: session.USER,
		IsArchive: false,
	}

	session, err := s.sessionRepo.CreateSession(ctx, s.sessionMapper.ToEntity(sessionInfo))
	if err != nil {
		return nil, custom_errors.ErrCreateUserSessionFailed
	}

	return s.sessionMapper.ToDomain(session), nil
}

func (s *SessionManager) GetListUserSession(ctx context.Context) ([]*session.SessionInfo, error) {
	sessions, err := s.sessionRepo.GetListSessionByOwnerType(ctx, session.USER.ToString())
	if err != nil {
		return nil, custom_errors.ErrGetUserSessionFailed
	}

	return s.sessionMapper.ToDomains(sessions), nil
}

func (s *SessionManager) GetUserSession(ctx context.Context, sessionID string) (*session.SessionInfo, error) {
	session, err := s.sessionRepo.GetUserSession(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to get user session", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetUserSessionFailed
	}

	if session == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	return s.sessionMapper.ToDomain(session), nil
}

func (s *SessionManager) GetUserChatHistory(ctx context.Context, sessionID string) (message.ListMessage, error) {
	session, err := s.sessionRepo.GetUserSessionWithMessages(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to get user chat history", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetUserSessionFailed
	}

	if session == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	domains, err := s.messageMapper.ToDomains(session.Messages)
	if err != nil {
		s.logger.Error("failed to get user chat history", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetChatHistoryFailed
	}

	return domains, nil
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
