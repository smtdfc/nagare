package chat

import (
	"context"

	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/core/session/manager"
	"github.com/smtdfc/nagare/shared/dtos/rest"
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

type Service struct {
	chatEventBus *EventBus
	sessionMgr   *manager.SessionManager
}

func (c *Service) SendMessage(ctx context.Context, request *rest.SendChatMessageRequest) error {
	_, err := c.sessionMgr.GetUserSession(ctx, request.SessionID)
	if err != nil {
		return err
	}

	c.chatEventBus.Publish(ctx, string(Topic), &SendMessageEvent{
		SessionID: request.SessionID,
		Text:      request.Text,
	})
	return nil
}

func (c *Service) CreateSession(ctx context.Context, request *rest.CreateChatSessionRequest) (*rest.CreateChatSessionResponse, error) {
	chatSession, err := c.sessionMgr.CreateUserSession(ctx, request.Title)
	if err != nil {
		return nil, err
	}

	return &rest.CreateChatSessionResponse{
		Session: toSessionDTO(chatSession),
	}, nil
}

func (c *Service) ListSessions(ctx context.Context) (*rest.ListChatSessionsResponse, error) {
	sessions, err := c.sessionMgr.GetListUserSession(ctx)
	if err != nil {
		return nil, err
	}

	return &rest.ListChatSessionsResponse{
		Sessions: helpers.SliceMap(sessions, toSessionDTO),
	}, nil
}

func (c *Service) GetHistory(ctx context.Context, sessionID string) (*rest.GetChatHistoryResponse, error) {
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
func NewService(sessionMgr *manager.SessionManager, chatEventBus *EventBus) *Service {
	return &Service{
		chatEventBus: chatEventBus,
		sessionMgr:   sessionMgr,
	}
}
