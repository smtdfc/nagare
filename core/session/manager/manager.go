package manager

import (
	"context"

	"github.com/google/uuid"
	config_mgr "github.com/smtdfc/nagare/core/config/manager"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type SessionManager struct {
	logger        *logger.BaseLogger
	pluginRepo    *repositories.PluginRepository
	sessionRepo   *repositories.SessionRepository
	messageRepo   *repositories.MessageRepository
	sessionMapper *mappers.SessionMapper
	messageMapper *mappers.MessageMapper
	configMgr     *config_mgr.ConfigManager
}

func (s *SessionManager) getDefaultLLMSettings(ctx context.Context) (uuid.UUID, string, error) {
	generalConf, err := s.configMgr.GetGeneralConfig(ctx)
	if err != nil {
		return uuid.Nil, "", err
	}

	if generalConf == nil || generalConf.DefaultLLMProvider == "" {
		return uuid.Nil, "", custom_errors.ErrMissingDefaultProvider
	}

	providerID, err := uuid.Parse(generalConf.DefaultLLMProvider)
	if err != nil {
		return uuid.Nil, "", err
	}

	return providerID, generalConf.DefaultLLMModel, nil
}

func newSessionState(sessionDomain *session.SessionInfo, messageDomains messages.ListMessage) *session.SessionState {
	return &session.SessionState{
		OwnerType:          sessionDomain.OwnerType,
		OwnerID:            sessionDomain.OwnerID,
		SessionID:          sessionDomain.ID,
		ChannelID:          sessionDomain.ChannelID,
		CurrentModel:       sessionDomain.CurrentLLMModel,
		CurrentLLMProvider: sessionDomain.LLMProviderID,
		Messages:           messageDomains,
	}
}

func (s *SessionManager) CreateUserSession(ctx context.Context, title string, ownerID string, llmProvider string, llmModel string) (*session.SessionInfo, error) {
	var llmProviderId = uuid.Nil
	var err error

	if llmProvider != "" {
		llmProviderId, err = uuid.Parse(llmProvider)
		if err != nil {
			return nil, custom_errors.ErrCreateSessionFailed
		}
	} else {
		llmProviderId, _, err = s.getDefaultLLMSettings(ctx)
		if err != nil {
			if err == custom_errors.ErrMissingDefaultProvider {
				return nil, err
			}

			s.logger.Error("Failed to get general config", "err", err)
			return nil, custom_errors.ErrCreateSessionFailed
		}
	}

	sessionInfo := &session.SessionInfo{
		Title:           title,
		OwnerID:         ownerID,
		OwnerType:       session.USER,
		IsArchive:       false,
		LLMProviderID:   llmProviderId,
		CurrentLLMModel: llmModel,
	}

	newSession, err := s.sessionRepo.Create(ctx, s.sessionMapper.ToEntity(sessionInfo))
	if err != nil {
		return nil, custom_errors.ErrCreateSessionFailed
	}

	return s.sessionMapper.ToDomain(newSession), nil
}

func (s *SessionManager) GetListUserSession(ctx context.Context, ownerID string) ([]*session.SessionInfo, error) {
	sessions, err := s.sessionRepo.FindByOwnerID(ctx, session.USER.ToString(), ownerID)
	if err != nil {
		return nil, custom_errors.ErrGetSessionFailed
	}

	return s.sessionMapper.ToDomains(sessions), nil
}

func (s *SessionManager) GetListUserSessionPage(ctx context.Context, ownerID string, offset int, limit int) ([]*session.SessionInfo, error) {
	sessions, err := s.sessionRepo.FindByOwnerIDPage(ctx, session.USER.ToString(), ownerID, offset, limit)
	if err != nil {
		return nil, custom_errors.ErrGetSessionFailed
	}

	return s.sessionMapper.ToDomains(sessions), nil
}

func (s *SessionManager) GetUserSession(ctx context.Context, sessionID string, ownerID string) (*session.SessionInfo, error) {
	userSession, err := s.sessionRepo.FindUserSession(ctx, sessionID, ownerID)
	if err != nil {
		s.logger.Error("failed to get session", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if userSession == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	return s.sessionMapper.ToDomain(userSession), nil
}

func (s *SessionManager) GetUserChatState(ctx context.Context, sessionID string, ownerID string) (*session.SessionState, error) {
	chatSession, err := s.sessionRepo.FindUserSessionWithMessages(ctx, sessionID, ownerID)
	if err != nil {
		s.logger.Error("failed to get chat history", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if chatSession == nil {
		return nil, custom_errors.ErrSessionNotFound
	}
	sessionDomain := s.sessionMapper.ToDomain(chatSession)
	domains, err := s.messageMapper.ToDomains(chatSession.Messages)
	if err != nil {
		s.logger.Error("failed to get chat history", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetChatHistoryFailed
	}

	return newSessionState(sessionDomain, domains), nil
}

func (s *SessionManager) GetUserChatStatePage(ctx context.Context, sessionID string, ownerID string, beforeID string, limit int) (*session.SessionState, error) {
	chatSession, err := s.sessionRepo.FindUserSession(ctx, sessionID, ownerID)
	if err != nil {
		s.logger.Error("failed to get chat history page", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if chatSession == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	domains, nextCursor, err := s.messageRepo.FindBySessionIDCursor(ctx, sessionID, beforeID, limit)
	if err != nil {
		s.logger.Error("failed to get chat history page", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetChatHistoryFailed
	}

	messageDomains, err := s.messageMapper.ToDomains(domains)
	if err != nil {
		s.logger.Error("failed to map chat history page", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetChatHistoryFailed
	}

	sessionDomain := s.sessionMapper.ToDomain(chatSession)
	state := newSessionState(sessionDomain, messageDomains)
	state.NextCursor = nextCursor

	return state, nil
}

func (s *SessionManager) PreparePluginSession(ctx context.Context, channelID string, targetID string) (*session.SessionInfo, error) {
	var err error
	plugin, err := s.pluginRepo.FindById(ctx, targetID)
	if err != nil {
		s.logger.Error("failed to prepare session", "channel_id", channelID, "target_id", targetID, "err", err)
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
		llmProviderID, defaultModel, err := s.getDefaultLLMSettings(ctx)
		if err != nil {
			if err == custom_errors.ErrMissingDefaultProvider {
				return nil, err
			}

			s.logger.Error("Failed to get general config", "err", err)
			return nil, custom_errors.ErrPreparePluginSessionFailed
		}

		sessionInfo := &session.SessionInfo{
			Title:           channelID,
			OwnerID:         plugin.ID.String(),
			OwnerType:       session.PLUGIN,
			IsArchive:       false,
			ChannelID:       channelID,
			LLMProviderID:   llmProviderID,
			CurrentLLMModel: defaultModel,
		}

		sessionEntity, err = s.sessionRepo.Create(ctx, s.sessionMapper.ToEntity(sessionInfo))
		if err != nil {
			return nil, custom_errors.ErrPreparePluginSessionFailed
		}
	}

	return s.sessionMapper.ToDomain(sessionEntity), nil
}

func (s *SessionManager) GetPluginSession(ctx context.Context, sessionID string, targetID string) (*session.SessionInfo, error) {
	plugin, err := s.pluginRepo.FindById(ctx, targetID)
	if err != nil {
		s.logger.Error("failed to get session", "session_id", sessionID, "target_id", targetID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}
	if plugin == nil {
		return nil, custom_errors.ErrPluginNotFound
	}

	sessionEntity, err := s.sessionRepo.FindPluginSession(ctx, sessionID, plugin.ID.String())
	if err != nil {
		s.logger.Error("failed to get session", "session_id", sessionID, "package_name", plugin.PackageName, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if sessionEntity == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	return s.sessionMapper.ToDomain(sessionEntity), nil
}

func (s *SessionManager) GetChatState(ctx context.Context, sessionID string) (*session.SessionState, error) {
	chatSession, err := s.sessionRepo.FindSessionWithMessages(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to get chat history", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if chatSession == nil {
		return nil, custom_errors.ErrSessionNotFound
	}

	sessionDomain := s.sessionMapper.ToDomain(chatSession)
	domains, err := s.messageMapper.ToDomains(chatSession.Messages)
	if err != nil {
		s.logger.Error("failed to get chat history", "session_id", sessionID, "err", err)
		return nil, custom_errors.ErrGetChatHistoryFailed
	}

	return newSessionState(sessionDomain, domains), nil
}

func (s *SessionManager) GetPluginChatState(ctx context.Context, sessionID string, pluginID string) (*session.SessionState, error) {
	plugin, err := s.pluginRepo.FindById(ctx, pluginID)
	if err != nil {
		s.logger.Error("failed to get session", "session_id", sessionID, "plugin_id", pluginID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if plugin == nil {
		return nil, custom_errors.ErrPluginNotFound
	}

	chatSession, err := s.sessionRepo.FindPluginSessionWithMessages(ctx, sessionID, plugin.ID.String())
	if err != nil {
		s.logger.Error("failed to get chat history", "session_id", sessionID, "plugin_id", pluginID, "err", err)
		return nil, custom_errors.ErrGetSessionFailed
	}

	if chatSession == nil {
		return nil, custom_errors.ErrSessionNotFound
	}
	sessionDomain := s.sessionMapper.ToDomain(chatSession)
	domains, err := s.messageMapper.ToDomains(chatSession.Messages)
	if err != nil {
		s.logger.Error("failed to get chat history", "session_id", sessionID, "plugin_id", pluginID, "err", err)
		return nil, custom_errors.ErrGetChatHistoryFailed
	}

	return newSessionState(sessionDomain, domains), nil
}

func (s *SessionManager) SaveHistory(ctx context.Context, sessionID string, pendingMessage messages.ListMessage) error {
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

func (s *SessionManager) ResetChatChannel(ctx context.Context, channelID string, pluginID string) error {
	plugin, err := s.pluginRepo.FindById(ctx, pluginID)
	if err != nil {
		s.logger.Error("failed to reset session", "plugin_id", pluginID, "err", err)
		return custom_errors.ErrPluginNotFound
	}

	if plugin == nil {
		return custom_errors.ErrPluginNotFound
	}

	chatSession, err := s.sessionRepo.FindByChannelID(ctx, session.PLUGIN.ToString(), plugin.ID.String(), channelID)
	if err != nil {
		s.logger.Error("failed to reset session", "channel_id", channelID, "err", err)
		return custom_errors.ErrResetSessionFailed
	}

	if chatSession == nil {
		return custom_errors.ErrSessionNotFound
	}

	err = s.messageRepo.DeleteBySessionID(ctx, chatSession.ID.String())
	if err != nil {
		s.logger.Error("failed to reset session", "session_id", chatSession.ID.String(), "err", err)
		return custom_errors.ErrResetSessionFailed
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
	configMgr *config_mgr.ConfigManager,
) *SessionManager {
	return &SessionManager{
		sessionRepo:   sessionRepo,
		sessionMapper: sessionMapper,
		messageRepo:   messageRepo,
		logger:        logger.With("module", "session-manager"),
		messageMapper: messageMapper,
		pluginRepo:    pluginRepo,
		configMgr:     configMgr,
	}
}
