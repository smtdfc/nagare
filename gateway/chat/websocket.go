package chat

import (
	"context"
	"fmt"

	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/core/session/manager"
	websocket2 "github.com/smtdfc/nagare/gateway/common/websocket"
	"github.com/smtdfc/nagare/shared/dtos/websocket"
)

type WebsocketHandler struct {
	sessionMgr *manager.SessionManager
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
func NewWebsocketHandler(sessionMgr *manager.SessionManager) *WebsocketHandler {
	return &WebsocketHandler{
		sessionMgr: sessionMgr,
	}
}
