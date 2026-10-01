package host

import (
	"context"

	"github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
)

func (h *PluginHost) handleHandshake(conn *PluginConnection, payload *websocket.Payload[any]) {
	data, err := getPayloadData[plugin.HandshakeEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin.HandshakeFailedEvent, &plugin.HandshakeFailedEventPayload{
			ID:    "",
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	h.logger.Info("Handshaking with plugin", "packageName", data.PackageName)
	pluginInfo, err := h.pluginMgr.ValidConnect(context.Background(), data.PackageName, data.ConnectCode)
	if err != nil {
		h.logger.Error("Handshake failed", "packageName", data.PackageName, "cause", err)
		_ = conn.Send(plugin.HandshakeFailedEvent, &plugin.HandshakeFailedEventPayload{
			ID:    data.ID,
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}

	scopes := make([]string, 0, len(pluginInfo.Features))
	for _, feature := range pluginInfo.Features {
		scopes = append(scopes, feature.ToString())
	}

	conn.pluginID = pluginInfo.ID.String()
	conn.pluginName = pluginInfo.PackageName
	conn.scopes = scopes

	h.mu.Lock()
	h.connections[conn.pluginID] = conn
	h.mu.Unlock()

	err = conn.Send(plugin.HandshakeSuccessEvent, &plugin.HandshakeSuccessEventPayload{
		ID: data.ID,
	}, payload.RequestID)
	if err != nil {
		return
	}
	h.logger.Info("Handshake success", "packageName", data.PackageName, "pluginID", conn.pluginID)
}
