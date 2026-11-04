package host

import (
	"context"

	"github.com/smtdfc/nagare/core/event_bus"
	core_plugin "github.com/smtdfc/nagare/core/plugin"
	tool_registry "github.com/smtdfc/nagare/core/tool/registry"
	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func (h *PluginHost) handleRegisterPluginTool(conn *PluginConnection, payload *contracts.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(contracts.RegisterPluginToolFailedEvent, &contracts.RegisterPluginToolFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[contracts.RegisterPluginToolEventPayload](payload)
	if err != nil {
		_ = conn.Send(contracts.RegisterPluginToolFailedEvent, &contracts.RegisterPluginToolFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !RequireScopes(conn.scopes, core_plugin.ToolFeature) {
		_ = conn.Send(contracts.RegisterPluginToolFailedEvent, &contracts.RegisterPluginToolFailedEventPayload{
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
		h.eventBus,
	)

	_ = conn.Send(contracts.RegisterPluginToolSuccessEvent, &contracts.RegisterPluginToolSuccessEventPayload{}, payload.RequestID)
}

func (h *PluginHost) handleRegisterToolCategories(conn *PluginConnection, payload *contracts.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(contracts.RegisterToolCategoriesFailedEvent, &contracts.RegisterToolCategoriesFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[contracts.RegisterToolCategoriesEventPayload](payload)
	if err != nil {
		_ = conn.Send(contracts.RegisterToolCategoriesFailedEvent, &contracts.RegisterToolCategoriesFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !RequireScopes(conn.scopes, core_plugin.ToolFeature) {
		_ = conn.Send(contracts.RegisterToolCategoriesFailedEvent, &contracts.RegisterToolCategoriesFailedEventPayload{
			Cause: "Plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	h.toolMgr.AddPluginToolCategories(conn.pluginID, data.Categories)
	h.logger.Info("Registered plugin tool categories", "categories", data.Categories, "pluginID", conn.pluginID)
	_ = conn.Send(contracts.RegisterToolCategoriesSuccessEvent, &contracts.RegisterToolCategoriesSuccessEventPayload{}, payload.RequestID)
}

func (h *PluginHost) handlePluginToolCallResult(conn *PluginConnection, payload *contracts.Payload[any]) {
	data, err := getPayloadData[contracts.PluginToolCallResultEventPayload](payload)
	if err != nil {
		return
	}

	h.eventBus.Publish(context.Background(), event_bus.PluginToolCallResultEvent, &event_bus.PluginToolCallResultEventPayload{
		RequestID: payload.RequestID,
		Result:    data.Result,
		Error:     data.Error,
	})
	h.logger.Debug("Plugin tool result received", "pluginID", conn.pluginID, "requestID", payload.RequestID)
}
