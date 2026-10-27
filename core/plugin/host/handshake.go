package host

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	plugin_dtos "github.com/smtdfc/nagare/plugin/dtos"
)

func (h *PluginHost) handleHandshake(conn *PluginConnection, payload *plugin_dtos.Payload[any]) {
	data, err := getPayloadData[plugin_dtos.HandshakeEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	h.logger.Info("Handshaking with plugin", "packageName", data.PackageName)
	pluginInfo, err := h.pluginMgr.ValidConnect(context.Background(), data.PackageName, data.ConnectCode)
	if err != nil {
		h.logger.Error("Handshake failed", "packageName", data.PackageName, "cause", err)
		_ = conn.Send(plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}

	groupName := getGroupName(pluginInfo.Features)

	if groupName != NormalPluginGroup {
		h.mu.Lock()

		if _, exist := h.specialGroups[groupName]; exist {
			h.mu.Unlock()
			h.logger.Error("Handshake blocked! The plugin belongs to a special category that has already been registered. It cannot be registered a second time for security reasons.", "packageName", data.PackageName)
			_ = conn.Send(plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
				Cause: custom_errors.ErrPluginBlockedBySpecialPluginGroupPolicy.Error(),
			}, payload.RequestID)
			conn.Close()
			return
		} else {
			h.specialGroups[groupName] = conn
			h.mu.Unlock()
		}
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

	err = conn.Send(plugin_dtos.HandshakeSuccessEvent, &plugin_dtos.HandshakeSuccessEventPayload{}, payload.RequestID)
	if err != nil {
		return
	}
	h.logger.Info("Handshake success", "packageName", data.PackageName, "pluginID", conn.pluginID)
}
