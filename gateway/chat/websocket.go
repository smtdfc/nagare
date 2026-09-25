package chat

import (
	"context"
	"fmt"

	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/session/manager"
	"github.com/smtdfc/nagare/dtos/websocket"
	"github.com/smtdfc/nagare/gateway/common/guards"
	websocket2 "github.com/smtdfc/nagare/gateway/common/websocket"
)

type ChatWebsocketHandler struct {
	sessionMgr *manager.SessionManager
	authGuard  *guards.AuthGuard
	logger     *logger.BaseLogger
}

func (c *ChatWebsocketHandler) OnAuth(s *melody.Session, w *websocket2.Coordinator, message *websocket.Payload[any]) {
	data, err := websocket2.GetData[websocket.AuthEventPayload](message)
	if err != nil {
		err := websocket2.SendMessage(s, websocket.AuthFailedEvent, &websocket.AuthFailedEventPayload{
			Cause: "failed to parsing payload",
		}, message.RequestID)
		if err != nil {
			return
		}

		return
	}

	if data == nil || data.Token == "" {
		err := websocket2.SendMessage(s, websocket.AuthFailedEvent, &websocket.AuthFailedEventPayload{
			Cause: "Unauthorized",
		}, message.RequestID)
		if err != nil {
			return
		}

		return
	}

	auth, err := c.authGuard.VerifyUserFromToken(data.Token)
	if err != nil {
		err := websocket2.SendMessage(s, websocket.AuthFailedEvent, &websocket.AuthFailedEventPayload{
			Cause: err.Error(),
		}, message.RequestID)
		if err != nil {
			return
		}

		return
	}

	s.Set("auth", &websocket2.AuthData{
		TargetType: "user",
		TargetID:   auth.ID,
		Scopes:     auth.Scopes,
	})

	err = websocket2.SendMessage(s, websocket.AuthSuccessEvent, &websocket.AuthSuccessEventPayload{}, message.RequestID)
	if err != nil {
		return
	}
}

func (c *ChatWebsocketHandler) OnListenMessage(s *melody.Session, w *websocket2.Coordinator, message *websocket.Payload[any]) {
	ctx := context.Background()
	data, err := websocket2.GetData[websocket.RegisterChatMessageListenerEventPayload](message)
	if err != nil {
		err := websocket2.SendMessage(s, websocket.RegisterChatListenerFailEvent, &websocket.RegisterChatListenerFailEventPayload{
			ID:    "",
			Cause: "failed to parsing payload",
		}, message.RequestID)
		if err != nil {
			return
		}
		return
	}

	if _, err := websocket2.GetAuth(s, "user"); err != nil {
		err := websocket2.SendMessage(s, websocket.RegisterChatListenerFailEvent, &websocket.RegisterChatListenerFailEventPayload{
			ID:    "",
			Cause: err.Error(),
		}, message.RequestID)
		if err != nil {
			return
		}
		return
	}

	auth, err := websocket2.GetAuth(s, "user")
	if err != nil {
		return
	}

	_, err = c.sessionMgr.GetUserSession(ctx, data.SessionID, auth.TargetID)
	if err != nil {
		err := websocket2.SendMessage(s, websocket.RegisterChatListenerFailEvent, &websocket.RegisterChatListenerFailEventPayload{
			ID:    data.ID,
			Cause: err.Error(),
		}, message.RequestID)
		if err != nil {
			return
		}

		return
	}

	w.JoinRoom(fmt.Sprintf("session:%s", data.SessionID), s)
	err = websocket2.SendMessage(s, websocket.RegisterChatListenerSuccessEvent, &websocket.RegisterChatListenerSuccessEventEventPayload{
		ID: data.ID,
	}, message.RequestID)
	if err != nil {
		return
	}
}

func (c *ChatWebsocketHandler) OnUnlistenMessage(s *melody.Session, w *websocket2.Coordinator, message *websocket.Payload[any]) {
	ctx := context.Background()
	data, err := websocket2.GetData[websocket.RegisterChatMessageListenerEventPayload](message)
	if err != nil {
		_ = websocket2.SendMessage(s, websocket.UnregisterChatListenerFailEvent, &websocket.RegisterChatListenerFailEventPayload{
			ID:    "",
			Cause: "failed to parsing payload",
		}, message.RequestID)
		return
	}

	auth, err := websocket2.GetAuth(s, "user")
	if err != nil {
		_ = websocket2.SendMessage(s, websocket.UnregisterChatListenerFailEvent, &websocket.RegisterChatListenerFailEventPayload{
			ID:    data.ID,
			Cause: err.Error(),
		}, message.RequestID)
		return
	}

	if _, err := c.sessionMgr.GetUserSession(ctx, data.SessionID, auth.TargetID); err != nil {
		_ = websocket2.SendMessage(s, websocket.UnregisterChatListenerFailEvent, &websocket.RegisterChatListenerFailEventPayload{
			ID:    data.ID,
			Cause: err.Error(),
		}, message.RequestID)
		return
	}

	w.LeaveRoom(fmt.Sprintf("session:%s", data.SessionID), s)
	_ = websocket2.SendMessage(s, websocket.UnregisterChatListenerSuccessEvent, &websocket.RegisterChatListenerSuccessEventEventPayload{
		ID: data.ID,
	}, message.RequestID)
}

// @Injectable
func NewWebsocketHandler(
	sessionMgr *manager.SessionManager,
	authGuard *guards.AuthGuard,
	logger *logger.BaseLogger,
) *ChatWebsocketHandler {
	return &ChatWebsocketHandler{
		sessionMgr: sessionMgr,
		authGuard:  authGuard,
		logger:     logger.With("module", "gateway:chat:websocket_handler"),
	}
}
