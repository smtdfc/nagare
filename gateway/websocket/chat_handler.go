package websocket

import (
	"context"
	"fmt"

	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/core/session/manager"
	"github.com/smtdfc/nagare/shared/dtos/websocket"
)

type ChatHandler struct {
	sessionMgr *manager.SessionManager
}

func (c *ChatHandler) OnListenMessage(s *melody.Session, w *Coordinator, message *websocket.WebsocketPayload[any]) {
	ctx := context.Background()
	data, err := GetData[websocket.RegisterChatMessageListenerEventPayload](message)
	if err != nil {
		err := SendMessage(s, websocket.RegisterChatListenerFailEvent, &websocket.RegisterChatListenerFailEventPayload{
			ID:    "",
			Cause: "failed to parsing payload",
		})
		if err != nil {
			return
		}
	}

	_, err = c.sessionMgr.GetUserSession(ctx, data.SessionID)
	if err != nil {
		err := SendMessage(s, websocket.RegisterChatListenerFailEvent, &websocket.RegisterChatListenerFailEventPayload{
			ID:    data.ID,
			Cause: err.Error(),
		})
		if err != nil {
			return
		}

		return
	}

	w.JoinRoom(fmt.Sprintf("session:%s", data.SessionID), s)
	err = SendMessage(s, websocket.RegisterChatListenerSuccessEvent, &websocket.RegisterChatListenerSuccessEventEventPayload{
		ID: data.ID,
	})
	if err != nil {
		return
	}
}

// @Injectable
func NewChatHandler(sessionMgr *manager.SessionManager) *ChatHandler {
	return &ChatHandler{
		sessionMgr: sessionMgr,
	}
}
