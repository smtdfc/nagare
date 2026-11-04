package host

import (
	"context"

	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func (h *PluginHost) handleHandshake(conn *PluginConnection, payload *contracts.Payload[any]) {
	data, err := getPayloadData[contracts.HandshakeEventPayload](payload)
	if err != nil {
		_ = conn.Send(contracts.HandshakeFailedEvent, &contracts.HandshakeFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	h.logger.Info("Handshaking with plugin", "packageName", data.PackageName)
	pluginInfo, err := h.pluginMgr.ValidConnect(context.Background(), data.PackageName, data.ConnectCode)
	if err != nil {
		h.logger.Error("Handshake failed", "packageName", data.PackageName, "cause", err)
		_ = conn.Send(contracts.HandshakeFailedEvent, &contracts.HandshakeFailedEventPayload{
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

	err = conn.Send(contracts.HandshakeSuccessEvent, &contracts.HandshakeSuccessEventPayload{}, payload.RequestID)
	if err != nil {
		return
	}
	h.logger.Info("Handshake success", "packageName", data.PackageName, "pluginID", conn.pluginID)
}
