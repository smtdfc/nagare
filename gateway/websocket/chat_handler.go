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

func (c *ChatHandler) OnListenMessage(s *melody.Session, w *WebsocketCoordinator, message *websocket.WebsocketPayload[any]) {
	ctx := context.Background()
	data, err := GetData[websocket.ChatListenMessageEvent](message)
	if err != nil {
		SendMessage(s, websocket.CHAT_LISTEN_MESSAGE_FAIL_EVENT, &websocket.ChatListenMessageFailEvent{
			ID:    "",
			Cause: "failed to parsing payload",
		})
	}

	_, err = c.sessionMgr.GetUserSession(ctx, data.SessionID)
	if err != nil {
		SendMessage(s, websocket.CHAT_LISTEN_MESSAGE_FAIL_EVENT, &websocket.ChatListenMessageFailEvent{
			ID:    data.ID,
			Cause: err.Error(),
		})

		return
	}

	w.JoinRoom(fmt.Sprintf("session:%s", data.SessionID), s)
	SendMessage(s, websocket.CHAT_LISTEN_MESSAGE_SUCCESS_EVENT, &websocket.ChatListenMessageSuccessEvent{
		ID: data.ID,
	})
}

// @Injectable
func NewChatHandler(sessionMgr *manager.SessionManager) *ChatHandler {
	return &ChatHandler{
		sessionMgr: sessionMgr,
	}
}
