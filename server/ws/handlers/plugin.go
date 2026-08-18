package handlers

import (
	"github.com/smtdfc/nagare/core/global"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/ws"
)

type PluginHandler struct {
}

func (h *PluginHandler) Register(i *ws.WsInstance, message *dto.WsMessage[any]) {
	payload, err := ws.GetPayload[dto.PluginRegisterEvent](message)
	if err != nil {
		ws.SendMessage(i, dto.WS_PLUGIN_REGISTER_FAILED, dto.PluginRegisterFailedEvent{
			Cause: "payload error: ",
		})
	}

	if !global.GlobalPluginMgr.HasConnectCode(payload.Code) {
		ws.SendMessage(i, dto.WS_PLUGIN_REGISTER_FAILED, dto.PluginRegisterFailedEvent{
			Cause: "invalid connect code",
		})
		return
	}

	i.Auth = &dto.AuthPayload{
		Target: "plugin",
	}

	ws.SendMessage(i, dto.WS_AUTH_SUCCESS, dto.PluginRegisterSuccessEvent{})
}

// @Injectable
func NewPluginHandler() *PluginHandler {
	return &PluginHandler{}
}
