package host

import (
	core_plugin "github.com/smtdfc/nagare/core/plugin"
	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func (h *PluginHost) handleInjectSystemPrompt(conn *PluginConnection, payload *contracts.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(contracts.InjectSystemPromptFailedEvent, &contracts.InjectSystemPromptFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[contracts.InjectSystemPromptEventPayload](payload)
	if err != nil {
		_ = conn.Send(contracts.InjectSystemPromptFailedEvent, &contracts.InjectSystemPromptFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !RequireScopes(conn.scopes, core_plugin.SystemPromptFeature) {
		_ = conn.Send(contracts.InjectSystemPromptFailedEvent, &contracts.InjectSystemPromptFailedEventPayload{
			Cause: "plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	if err := h.resourceMgr.SetSystemPrompt(conn.pluginID, data.Content); err != nil {
		_ = conn.Send(contracts.InjectSystemPromptFailedEvent, &contracts.InjectSystemPromptFailedEventPayload{
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}

	_ = conn.Send(contracts.InjectSystemPromptSuccessEvent, &contracts.InjectSystemPromptSuccessEventPayload{}, payload.RequestID)
}

func (h *PluginHost) handleClearSystemPrompt(conn *PluginConnection, payload *contracts.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(contracts.ClearSystemPromptFailedEvent, &contracts.ClearSystemPromptFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	if !RequireScopes(conn.scopes, core_plugin.SystemPromptFeature) {
		_ = conn.Send(contracts.ClearSystemPromptFailedEvent, &contracts.ClearSystemPromptFailedEventPayload{
			Cause: "plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	h.resourceMgr.ClearSystemPrompt(conn.pluginID)

	_ = conn.Send(contracts.ClearSystemPromptSuccessEvent, &contracts.ClearSystemPromptSuccessEventPayload{}, payload.RequestID)
}
