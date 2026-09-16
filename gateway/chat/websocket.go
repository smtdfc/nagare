package chat

import (
	"context"
	"fmt"

	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/core/session/manager"
	"github.com/smtdfc/nagare/dtos/websocket"
	"github.com/smtdfc/nagare/gateway/common/guards"
	websocket2 "github.com/smtdfc/nagare/gateway/common/websocket"
)

type WebsocketHandler struct {
	sessionMgr *manager.SessionManager
	authGuard  *guards.AuthGuard
}

func (c *WebsocketHandler) OnAuth(s *melody.Session, w *websocket2.Coordinator, message *websocket.Payload[any]) {
	data, err := websocket2.GetData[websocket.AuthEventPayload](message)
	if err != nil {
		err := websocket2.SendMessage(s, websocket.AuthFailedEvent, &websocket.AuthFailedEventPayload{
			Cause: "failed to parsing payload",
		}, message.RequestID)
		if err != nil {
			return
		}
	}

	if data != nil || data.Token == "" {
		err := websocket2.SendMessage(s, websocket.AuthFailedEvent, &websocket.AuthFailedEventPayload{
			Cause: "Unauthorized",
		}, message.RequestID)
		if err != nil {
			return
		}
	}

	auth, err := c.authGuard.VerifyUserFromToken(data.Token)
	if err != nil {
		err := websocket2.SendMessage(s, websocket.AuthFailedEvent, &websocket.AuthFailedEventPayload{
			Cause: err.Error(),
		}, message.RequestID)
		if err != nil {
			return
		}
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

func (c *WebsocketHandler) OnListenMessage(s *melody.Session, w *websocket2.Coordinator, message *websocket.Payload[any]) {
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

	_, err = c.sessionMgr.GetUserSession(ctx, data.SessionID)
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

// @Injectable
func NewWebsocketHandler(
	sessionMgr *manager.SessionManager,
	authGuard *guards.AuthGuard,
) *WebsocketHandler {
	return &WebsocketHandler{
		sessionMgr: sessionMgr,
		authGuard:  authGuard,
	}
}
