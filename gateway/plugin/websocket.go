package plugin

import (
	"context"

	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/core/plugin/manager"
	websocket2 "github.com/smtdfc/nagare/gateway/common/websocket"
	plugin_dtos "github.com/smtdfc/nagare/shared/dtos/plugin"
	"github.com/smtdfc/nagare/shared/dtos/websocket"
)

type WebsocketHandler struct {
	pluginMgr *manager.PluginManager
}

func (w *WebsocketHandler) OnHandshakeEvent(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
	ctx := context.Background()
	data, err := websocket2.GetData[plugin_dtos.HandshakeEventPayload](message)
	if err != nil {
		err := websocket2.SendMessage(s, plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			ID:    "",
			Cause: "failed to parsing payload",
		})
		if err != nil {
			return
		}
	}

	err = w.pluginMgr.ValidConnect(ctx, data.PluginID, data.ConnectCode)
	if err != nil {
		err := websocket2.SendMessage(s, plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			ID:    data.ID,
			Cause: err.Error(),
		})
		if err != nil {
			return
		}
	}

	s.Set("auth", websocket2.AuthData{
		TargetType: "plugin",
		TargetID:   data.PluginID,
	})

	err = websocket2.SendMessage(s, plugin_dtos.HandshakeSuccessEvent, &plugin_dtos.HandshakeSuccessEventPayload{
		ID: data.ID,
	})
	if err != nil {
		return
	}
}

// @Injectable
func NewWebsocketHandler(
	pluginMgr *manager.PluginManager,
) *WebsocketHandler {
	return &WebsocketHandler{
		pluginMgr: pluginMgr,
	}
}
