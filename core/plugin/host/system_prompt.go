package host

import (
	core_plugin "github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/prompt"
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

	if len(data.Content) > contracts.MaxSystemPromptBytes {
		_ = conn.Send(contracts.InjectSystemPromptFailedEvent, &contracts.InjectSystemPromptFailedEventPayload{
			Cause: "system prompt too large",
		}, payload.RequestID)
		return
	}

	prompt.SetPluginSystemPrompt(conn.pluginID, data.Content)
	h.logger.Info(
		"Injected plugin system prompt",
		"pluginID", conn.pluginID,
		"bytes", len(data.Content),
	)

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

	prompt.ClearPluginSystemPrompt(conn.pluginID)
	h.logger.Info(
		"Cleared plugin system prompt",
		"pluginID", conn.pluginID,
	)

	_ = conn.Send(contracts.ClearSystemPromptSuccessEvent, &contracts.ClearSystemPromptSuccessEventPayload{}, payload.RequestID)
}
