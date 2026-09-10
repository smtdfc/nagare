package services

import (
	"context"

	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/core/session/manager"
	"github.com/smtdfc/nagare/gateway/event"
	"github.com/smtdfc/nagare/shared/dtos/rest"
	"github.com/smtdfc/nagare/shared/event_bus"
	"github.com/smtdfc/nagare/shared/helpers"
)

func toSessionDTO(s *session.SessionInfo) *rest.Session {
	if s == nil {
		return nil
	}
	return &rest.Session{
		ID:    s.ID,
		Title: s.Title,
	}
}

type ChatService struct {
	chatEventBus *event_bus.EventBus[event.ChatEventPayload]
	sessionMgr   *manager.SessionManager
}

func (c *ChatService) SendMessage(ctx context.Context, request *rest.SendChatMessageRequest) error {
	_, err := c.sessionMgr.GetUserSession(ctx, request.SessionID)
	if err != nil {
		return err
	}

	c.chatEventBus.Publish(ctx, string(event.ChatTopic), &event.ChatSendMessageEvent{
		SessionID: request.SessionID,
		Text:      request.Text,
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
	messages, err := c.sessionMgr.GetUserChatHistory(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &rest.GetChatHistoryResponse{
		SessionID: sessionID,
		Messages:  messages,
	}, nil
}

// @Injectable
func NewChatService(sessionMgr *manager.SessionManager, busSys *event.AppEventBusSystem) *ChatService {
	return &ChatService{
		chatEventBus: busSys.ChatEventBus,
		sessionMgr:   sessionMgr,
	}
}
