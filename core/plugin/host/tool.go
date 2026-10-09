package host

import (
	"context"
	"slices"

	"github.com/smtdfc/nagare/core/event_bus"
	core_plugin "github.com/smtdfc/nagare/core/plugin"
	tool_registry "github.com/smtdfc/nagare/core/tool/registry"
	"github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
)

func (h *PluginHost) handleRegisterPluginTool(conn *PluginConnection, payload *websocket.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin.RegisterPluginToolFailedEvent, &plugin.RegisterPluginToolFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin.RegisterPluginToolEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin.RegisterPluginToolFailedEvent, &plugin.RegisterPluginToolFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !slices.Contains(conn.scopes, core_plugin.ToolFeature.ToString()) {
		_ = conn.Send(plugin.RegisterPluginToolFailedEvent, &plugin.RegisterPluginToolFailedEventPayload{
			Cause: "Plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	h.logger.Info("Registering plugin tool", "name", data.Tool.Name, "pluginID", conn.pluginID)
	tool_registry.RegisterPluginTool(
		data.Tool.Name,
		data.Tool.Args,
		data.Tool.Description,
		data.Tool.Categories,
		conn.pluginID,
		true,
		h.chatEventBus,
	)

	_ = conn.Send(plugin.RegisterPluginToolSuccessEvent, &plugin.RegisterPluginToolSuccessEventPayload{}, payload.RequestID)
}

func (h *PluginHost) handleRegisterToolCategories(conn *PluginConnection, payload *websocket.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin.RegisterToolCategoriesFailedEvent, &plugin.RegisterToolCategoriesFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin.RegisterToolCategoriesEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin.RegisterToolCategoriesFailedEvent, &plugin.RegisterToolCategoriesFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !slices.Contains(conn.scopes, core_plugin.ToolFeature.ToString()) {
		_ = conn.Send(plugin.RegisterToolCategoriesFailedEvent, &plugin.RegisterToolCategoriesFailedEventPayload{
			Cause: "Plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	h.toolMgr.AddPluginToolCategories(conn.pluginID, data.Categories)
	h.logger.Info("Registered plugin tool categories", "categories", data.Categories, "pluginID", conn.pluginID)
	_ = conn.Send(plugin.RegisterToolCategoriesSuccessEvent, &plugin.RegisterToolCategoriesSuccessEventPayload{}, payload.RequestID)
}

func (h *PluginHost) handlePluginToolCallResult(conn *PluginConnection, payload *websocket.Payload[any]) {
	data, err := getPayloadData[plugin.PluginToolCallResultEventPayload](payload)
	if err != nil {
		return
	}

	h.chatEventBus.Publish(context.Background(), event_bus.PluginToolCallResultEvent, &event_bus.PluginToolCallResultEventPayload{
		RequestID: payload.RequestID,
		Result:    data.Result,
		Error:     data.Error,
	})
	h.logger.Debug("Plugin tool result received", "pluginID", conn.pluginID, "requestID", payload.RequestID)
}
