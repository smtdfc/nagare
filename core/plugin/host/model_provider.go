package host

import (
	"context"
	"fmt"

	core_plugin "github.com/smtdfc/nagare/core/plugin"
	plugin_dtos "github.com/smtdfc/nagare/plugin/dtos"
)

func (h *PluginHost) handleGetModelProviders(conn *PluginConnection, payload *plugin_dtos.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin_dtos.GetModelProvidersFailedEvent, &plugin_dtos.GetModelProvidersFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}
	fmt.Println(2)
	if !RequireScopes(conn.scopes, core_plugin.ModelProviderFeature) {
		fmt.Println(2)
		_ = conn.Send(plugin_dtos.GetModelProvidersFailedEvent, &plugin_dtos.GetModelProvidersFailedEventPayload{
			Cause: "plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	ctx := context.Background()
	providers, err := h.modelProvider.GetAllProvider(ctx)
	fmt.Println(3)
	if err != nil {
		_ = conn.Send(plugin_dtos.GetModelProvidersFailedEvent, &plugin_dtos.GetModelProvidersFailedEventPayload{
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}
	fmt.Println(providers)
	providerMap := make(map[string]string)
	for _, provider := range providers {
		providerMap[provider.ID.String()] = provider.Name
	}

	_ = conn.Send(plugin_dtos.GetModelProvidersSuccessEvent, &plugin_dtos.GetModelProvidersSuccessEventPayload{
		List: providerMap,
	}, payload.RequestID)
}
