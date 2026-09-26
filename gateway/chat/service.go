package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/core/session/manager"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/shared/helpers"
)

func toSessionDTO(s *session.SessionInfo) *rest.Session {
	if s == nil {
		return nil
	}

	return &rest.Session{
		ID:    s.ID.String(),
		Title: s.Title,
	}
}

type ChatService struct {
	chatEventBus *event_bus.CoreEventBus
	sessionMgr   *manager.SessionManager
	logger       *logger.BaseLogger
}

func (c *ChatService) SendMessage(ctx context.Context, ownerID string, request *rest.SendChatMessageRequest) error {
	_, err := c.sessionMgr.GetUserSession(ctx, request.SessionID, ownerID)
	if err != nil {
		return err
	}

	c.chatEventBus.Publish(ctx, event_bus.SendEvent, &event_bus.SendMessageEventPayload{
		RequestID:  uuid.New().String(),
		SessionID:  request.SessionID,
		Text:       request.Text,
		SenderType: event_bus.User,
		SenderID:   ownerID,
	})
	return nil
}

func (c *ChatService) CreateSession(ctx context.Context, ownerID string, request *rest.CreateChatSessionRequest) (*rest.CreateChatSessionResponse, error) {
	chatSession, err := c.sessionMgr.CreateUserSession(ctx, request.Title, ownerID)
	if err != nil {
		return nil, err
	}

	return &rest.CreateChatSessionResponse{
		Session: toSessionDTO(chatSession),
	}, nil
}

func (c *ChatService) ListSessions(ctx context.Context, ownerID string, offset int, limit int) (*rest.ListChatSessionsResponse, error) {
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
		Sessions: helpers.SliceMap(sessions, toSessionDTO),
	}, nil
}

func (c *ChatService) GetHistory(ctx context.Context, ownerID string, sessionID string) (*rest.GetChatHistoryResponse, error) {
	sessionHistory, err := c.sessionMgr.GetUserChatHistory(ctx, sessionID, ownerID)
	if err != nil {
		return nil, err
	}

	return &rest.GetChatHistoryResponse{
		SessionID: sessionID,
		Messages:  sessionHistory.Messages,
	}, nil
}

func (c *ChatService) GetHistoryPage(ctx context.Context, ownerID string, sessionID string, beforeID string, limit int) (*rest.GetChatHistoryResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	sessionHistory, err := c.sessionMgr.GetUserChatHistoryPage(ctx, sessionID, ownerID, beforeID, limit)
	if err != nil {
		return nil, err
	}

	return &rest.GetChatHistoryResponse{
		SessionID:  sessionID,
		Messages:   sessionHistory.Messages,
		NextCursor: sessionHistory.NextCursor,
	}, nil
}
func (c *ChatService) GetSession(ctx context.Context, ownerID string, sessionID string) (*rest.GetChatSessionResponse, error) {
	s, err := c.sessionMgr.GetUserSession(ctx, sessionID, ownerID)
	if err != nil {
		return nil, err
	}

	return &rest.GetChatSessionResponse{
		Session: toSessionDTO(s),
	}, nil
}

// @Injectable
func NewService(sessionMgr *manager.SessionManager, chatEventBus *event_bus.CoreEventBus, logger *logger.BaseLogger) *ChatService {
	return &ChatService{
		chatEventBus: chatEventBus,
		sessionMgr:   sessionMgr,
		logger:       logger.With("module", "gateway:chat:service"),
	}
}
