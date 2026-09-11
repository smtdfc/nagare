package plugin

import (
	"context"

	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/plugin/manager"
	websocket2 "github.com/smtdfc/nagare/gateway/common/websocket"
	plugin_dtos "github.com/smtdfc/nagare/shared/dtos/plugin"
	"github.com/smtdfc/nagare/shared/dtos/websocket"
)

type WebsocketHandler struct {
	pluginMgr *manager.PluginManager
	logger    *logger.BaseLogger
}

func (w *WebsocketHandler) OnHandshakeEvent(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {

	ctx := context.Background()
	data, err := websocket2.GetData[plugin_dtos.HandshakeEventPayload](message)
	if err != nil {
		err := websocket2.SendMessage(s, plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			ID:    "",
			Cause: "failed to parsing payload",
		}, message.RequestID)
		if err != nil {
			return
		}
	}

	w.logger.Info("Handshaking with plugin", "pluginID", data.PluginID)
	err = w.pluginMgr.ValidConnect(ctx, data.PluginID, data.ConnectCode)
	if err != nil {
		w.logger.Error("Handshake failed", "pluginID", data.PluginID, "cause", err)
		err := websocket2.SendMessage(s, plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			ID:    data.ID,
			Cause: err.Error(),
		}, message.RequestID)
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
	}, message.RequestID)
	if err != nil {
		return
	}
	w.logger.Info("Handshake success", "pluginID", data.PluginID)
}

// @Injectable
func NewWebsocketHandler(
	pluginMgr *manager.PluginManager,
	logger *logger.BaseLogger,
) *WebsocketHandler {
	return &WebsocketHandler{
		pluginMgr: pluginMgr,
		logger:    logger,
	}
}
