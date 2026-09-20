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

func (c *ChatService) SendMessage(ctx context.Context, request *rest.SendChatMessageRequest) error {
	_, err := c.sessionMgr.GetUserSession(ctx, request.SessionID)
	if err != nil {
		return err
	}

	c.chatEventBus.Publish(ctx, event_bus.SendEvent, &event_bus.SendMessageEventPayload{
		RequestID:  uuid.New().String(),
		SessionID:  request.SessionID,
		Text:       request.Text,
		SenderType: event_bus.User,
		SenderID:   "",
	})
	return nil
}

func (c *ChatService) CreateSession(ctx context.Context, request *rest.CreateChatSessionRequest) (*rest.CreateChatSessionResponse, error) {
	chatSession, err := c.sessionMgr.CreateUserSession(ctx, request.Title)
	if err != nil {
		return nil, err
	}

	return &rest.CreateChatSessionResponse{
		Session: toSessionDTO(chatSession),
	}, nil
}

func (c *ChatService) ListSessions(ctx context.Context) (*rest.ListChatSessionsResponse, error) {
	sessions, err := c.sessionMgr.GetListUserSession(ctx)
	if err != nil {
		return nil, err
	}

	return &rest.ListChatSessionsResponse{
		Sessions: helpers.SliceMap(sessions, toSessionDTO),
	}, nil
}

func (c *ChatService) GetHistory(ctx context.Context, sessionID string) (*rest.GetChatHistoryResponse, error) {
	sessionHistory, err := c.sessionMgr.GetUserChatHistory(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &rest.GetChatHistoryResponse{
		SessionID: sessionID,
		Messages:  sessionHistory.Messages,
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
