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
	pluginRepo    *repositories.PluginRepository
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
		return nil, custom_errors.ErrCreateSessionFailed
	}

	return s.sessionMapper.ToDomain(newSession), nil
}

func (s *SessionManager) GetListUserSession(ctx context.Context) ([]*session.SessionInfo, error) {
	sessions, err := s.sessionRepo.FindByOwnerType(ctx, session.USER.ToString())
	if err != nil {
		return nil, custom_errors.ErrGetSessionFailed
	}

	return s.sessionMapper.ToDomains(sessions), nil
}

func (s *SessionManager) GetUserSession(ctx context.Context, sessionID string) (*session.SessionInfo, error) {
	userSession, err := s.sessionRepo.FindUserSession(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to get session", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if userSession == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	return s.sessionMapper.ToDomain(userSession), nil
}

func (s *SessionManager) GetUserChatHistory(ctx context.Context, sessionID string) (message.ListMessage, error) {
	chatSession, err := s.sessionRepo.FindUserSessionWithMessages(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to get chat history", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if chatSession == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	domains, err := s.messageMapper.ToDomains(chatSession.Messages)
	if err != nil {
		s.logger.Error("failed to get chat history", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetChatHistoryFailed
	}

	return domains, nil
}

func (s *SessionManager) PreparePluginSession(ctx context.Context, channelID string, pluginID string) (*session.SessionInfo, error) {
	var err error
	plugin, err := s.pluginRepo.FindByPluginId(ctx, pluginID)
	if err != nil {
		s.logger.Error("failed to prepare session", "channel_id", channelID, "plugin_id", pluginID, "err", err)
		return nil, custom_errors.ErrPreparePluginSessionFailed
	}
	if plugin == nil {
		return nil, custom_errors.ErrPluginNotFound
	}

	sessionEntity, err := s.sessionRepo.FindByChannelID(
		ctx,
		session.PLUGIN.ToString(),
		plugin.ID.String(),
		channelID,
	)
	if err != nil {
		return nil, custom_errors.ErrPreparePluginSessionFailed
	}

	if sessionEntity == nil {
		sessionInfo := &session.SessionInfo{
			Title:     channelID,
			OwnerID:   plugin.ID,
			OwnerType: session.PLUGIN,
			IsArchive: false,
			ChannelID: channelID,
		}

		sessionEntity, err = s.sessionRepo.Create(ctx, s.sessionMapper.ToEntity(sessionInfo))
		if err != nil {
			return nil, custom_errors.ErrPreparePluginSessionFailed
		}
	}

	return s.sessionMapper.ToDomain(sessionEntity), nil
}

func (s *SessionManager) GetPluginSession(ctx context.Context, sessionID string, pluginID string) (*session.SessionInfo, error) {
	plugin, err := s.pluginRepo.FindByPluginId(ctx, pluginID)
	if err != nil {
		s.logger.Error("failed to get session", "session_id", sessionID, "plugin_id", pluginID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}
	if plugin == nil {
		return nil, custom_errors.ErrPluginNotFound
	}

	sessionEntity, err := s.sessionRepo.FindPluginSession(ctx, sessionID, plugin.ID.String())
	if err != nil {
		s.logger.Error("failed to get session", "session_id", sessionID, "plugin_id", pluginID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if sessionEntity == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	return s.sessionMapper.ToDomain(sessionEntity), nil
}

func (s *SessionManager) GetPluginChatHistory(ctx context.Context, sessionID string, pluginID string) (message.ListMessage, error) {
	plugin, err := s.pluginRepo.FindByPluginId(ctx, pluginID)
	if err != nil {
		s.logger.Error("failed to get session", "session_id", sessionID, "plugin_id", pluginID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	chatSession, err := s.sessionRepo.FindPluginSessionWithMessages(ctx, sessionID, plugin.ID.String())
	if err != nil {
		s.logger.Error("failed to get chat history", "session_id", sessionID, "plugin_id", pluginID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if chatSession == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	domains, err := s.messageMapper.ToDomains(chatSession.Messages)
	if err != nil {
		s.logger.Error("failed to get chat history", "session_id", sessionID, "plugin_id", pluginID, "err", err)
		return nil, custom_errors.ErrGetChatHistoryFailed
	}

	return domains, nil
}

func (s *SessionManager) SaveHistory(ctx context.Context, sessionID string, pendingMessage message.ListMessage) error {
	chatSession, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to save chat history", "session_id", sessionID, "err", err)
		return custom_errors.ErrSaveSessionFailed
	}

	if chatSession == nil {
		return custom_errors.ErrSessionNotFound
	}

	entities, err := s.messageMapper.ToEntities(pendingMessage, sessionID)
	if err != nil {
		s.logger.Error("failed to save chat history", "session_id", sessionID, "err", err)
		return custom_errors.ErrSaveSessionFailed
	}

	err = s.messageRepo.CreateBatch(ctx, entities, 200)
	if err != nil {
		s.logger.Error("failed to save chat history", "session_id", sessionID, "err", err)
		return custom_errors.ErrSaveSessionFailed
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
	pluginRepo *repositories.PluginRepository,
) *SessionManager {
	return &SessionManager{
		sessionRepo:   sessionRepo,
		sessionMapper: sessionMapper,
		messageRepo:   messageRepo,
		logger:        logger,
		messageMapper: messageMapper,
		pluginRepo:    pluginRepo,
	}
}
