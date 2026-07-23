package handlers

import (
	"github.com/smtdfc/nagare/core/providers"
	"github.com/smtdfc/nagare/server/ws"
	"github.com/smtdfc/nagare/shared/dto"
)

type ChatHandler struct{}

func (h *ChatHandler) CreateSession(i *ws.WsInstance) {
	sessionID, err := providers.GlobalSessionManager.CreateSession()
	if err != nil {
		ws.SendMessage(i, dto.WS_CREATE_SESSION_FAILED, dto.CreateSessionFailed{
			Cause: err.Error(),
		})
	}

	ws.SendMessage(i, dto.WS_CREATE_SESSION_SUCCESS, dto.CreateSessionSuccess{
		ID: sessionID,
	})
}

// @Injectable
func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}
