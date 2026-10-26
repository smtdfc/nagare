package chat

import (
	"context"

	"github.com/google/uuid"
	config_mgr "github.com/smtdfc/nagare/core/config/manager"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/session/manager"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type Service struct {
	configMgr     *config_mgr.ConfigManager
	chatEventBus  *event_bus.CoreEventBus
	sessionMgr    *manager.SessionManager
	logger        *logger.BaseLogger
	sessionMapper ISessionDtoMapper
}

func (c *Service) SendMessage(ctx context.Context, ownerID string, request *rest.SendChatMessageRequest) (*rest.SendChatMessageResponse, error) {
	invokeID := uuid.New().String()
	_, err := c.sessionMgr.GetUserSession(ctx, request.SessionID, ownerID)
	if err != nil {
		return nil, err
	}

	c.chatEventBus.Publish(ctx, event_bus.SendChatMessageEvent, &event_bus.SendChatMessageEventPayload{
		InvokeID:  invokeID,
		RequestID: uuid.New().String(),
		SessionID: request.SessionID,
		Messages: messages.ListMessage{
			messages.NewTextMessage(messages.USER, request.Text),
		},
		SenderType: event_bus.User,
		SenderID:   ownerID,
	})

	return &rest.SendChatMessageResponse{
		InvokeID: invokeID,
	}, nil
}

func (c *Service) CancelChat(ctx context.Context, ownerID string, request *rest.CancelChatRequest) error {
	_, err := c.sessionMgr.GetUserSession(ctx, request.SessionID, ownerID)
	if err != nil {
		return err
	}

	c.chatEventBus.Publish(ctx, event_bus.CancelChatEvent, &event_bus.CancelChatEventPayload{
		InvokeID:  request.InvokeID,
		RequestID: uuid.New().String(),
	})

	return nil
}

func (c *Service) CreateSession(ctx context.Context, ownerID string, request *rest.CreateChatSessionRequest) (*rest.CreateChatSessionResponse, error) {
	chatSession, err := c.sessionMgr.CreateUserSession(ctx, request.Title, ownerID, request.CurrentModelProvider, request.CurrentModel)
	if err != nil {
		return nil, err
	}

	return &rest.CreateChatSessionResponse{
		Session: c.sessionMapper.ToSessionDTO(chatSession),
	}, nil
}

func (c *Service) ListSessions(ctx context.Context, ownerID string, offset int, limit int) (*rest.ListChatSessionsResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	sessions, err := c.sessionMgr.GetListUserSessionPage(ctx, ownerID, offset, limit)
	if err != nil {
		return nil, err
	}

	return &rest.ListChatSessionsResponse{
		Sessions: c.sessionMapper.ToSessionDTOs(sessions),
	}, nil
}

// func (c *ChatService) GetHistory(ctx context.Context, ownerID string, sessionID string) (*rest.GetChatHistoryResponse, error) {
// 	sessionHistory, err := c.sessionMgr.GetUserChatState(ctx, sessionID, ownerID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &rest.GetChatHistoryResponse{
// 		SessionID: sessionID,
// 		Messages:  sessionHistory.Messages,
// 	}, nil
// }

func (c *Service) GetHistoryPage(ctx context.Context, ownerID string, sessionID string, beforeID string, limit int) (*rest.GetChatHistoryResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	sessionHistory, err := c.sessionMgr.GetUserChatStatePage(ctx, sessionID, ownerID, beforeID, limit)
	if err != nil {
		return nil, err
	}

	return &rest.GetChatHistoryResponse{
		SessionID:  sessionID,
		Messages:   sessionHistory.Messages,
		NextCursor: sessionHistory.NextCursor,
	}, nil
}
func (c *Service) GetSession(ctx context.Context, ownerID string, sessionID string) (*rest.GetChatSessionResponse, error) {
	s, err := c.sessionMgr.GetUserSession(ctx, sessionID, ownerID)
	if err != nil {
		return nil, err
	}

	return &rest.GetChatSessionResponse{
		Session: c.sessionMapper.ToSessionDTO(s),
	}, nil
}

func (c *Service) UpdateModelSettings(ctx context.Context, ownerID string, sessionID string, request *rest.UpdateChatSessionModelSettingsRequest) (*rest.UpdateChatSessionModelSettingsResponse, error) {
	chatSession, err := c.sessionMgr.UpdateUserSessionModelSettings(
		ctx,
		sessionID,
		ownerID,
		request.CurrentModelProvider,
		request.CurrentModel,
	)
	if err != nil {
		return nil, err
	}

	return &rest.UpdateChatSessionModelSettingsResponse{
		Session: c.sessionMapper.ToSessionDTO(chatSession),
	}, nil
}

func (c *Service) DeleteSession(ctx context.Context, ownerID string, sessionID string) error {
	return c.sessionMgr.DeleteUserSession(ctx, sessionID, ownerID)
}

func (c *Service) ArchiveSession(ctx context.Context, ownerID string, sessionID string, request *rest.ArchiveChatSessionRequest) (*rest.ArchiveChatSessionResponse, error) {
	chatSession, err := c.sessionMgr.ArchiveUserSession(ctx, sessionID, ownerID, request.IsArchive)
	if err != nil {
		return nil, err
	}

	return &rest.ArchiveChatSessionResponse{
		Session: c.sessionMapper.ToSessionDTO(chatSession),
	}, nil
}

func (c *Service) DuplicateSession(ctx context.Context, ownerID string, sessionID string) (*rest.DuplicateChatSessionResponse, error) {
	chatSession, err := c.sessionMgr.DuplicateUserSession(ctx, sessionID, ownerID)
	if err != nil {
		return nil, err
	}

	return &rest.DuplicateChatSessionResponse{
		Session: c.sessionMapper.ToSessionDTO(chatSession),
	}, nil
}

// @Injectable
func NewService(sessionMgr *manager.SessionManager, configMgr *config_mgr.ConfigManager, chatEventBus *event_bus.CoreEventBus, logger *logger.BaseLogger, sessionMapper ISessionDtoMapper) *Service {
	return &Service{
		configMgr:     configMgr,
		chatEventBus:  chatEventBus,
		sessionMgr:    sessionMgr,
		logger:        logger.With("module", "gateway:chat:service"),
		sessionMapper: sessionMapper,
	}
}
