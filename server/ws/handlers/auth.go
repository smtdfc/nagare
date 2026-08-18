package handlers

import (
	"fmt"

	"github.com/smtdfc/nagare/server/config"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/ws"
)

type AuthHandler struct {
	conf *config.ServerConfig
}

func (h *AuthHandler) Auth(i *ws.WsInstance, message *dto.WsMessage[any]) {
	payload, err := ws.GetPayload[dto.AuthRequestEvent](message)
	if err != nil {
		ws.SendMessage(i, dto.WS_AUTH_FAILED, dto.WsAuthFailedEvent{

			Cause: "payload error: ",
		})
	}

	token := payload.Token
	authPayload, err := helpers.VerifyToken(h.conf.PublicKey, token)
	if err != nil {
		fmt.Println("err", err)
		ws.SendMessage(i, dto.WS_AUTH_FAILED, dto.WsAuthFailedEvent{

			Cause: "unauthorized",
		})
	}

	i.Auth = authPayload

	ws.SendMessage(i, dto.WS_AUTH_SUCCESS, dto.WsAuthSuccessEvent{})
}

// @Injectable
func NewAuthHandler(conf *config.ServerConfig) *AuthHandler {
	return &AuthHandler{
		conf: conf,
	}
}
