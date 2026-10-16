package client

import (
	"context"
	"errors"

	plugin_dtos "github.com/smtdfc/nagare/plugin/dtos"
)

func (p *PluginClient) GetListModelProvider(ctx context.Context) (map[string]string, error) {
	payload := plugin_dtos.GetModelProvidersEventPayload{}
	var providerMap = make(map[string]string)
	resp, err := sendAndWait[plugin_dtos.GetModelProvidersSuccessEventPayload, plugin_dtos.GetModelProvidersFailedEventPayload](
		p,
		ctx,
		plugin_dtos.GetModelProvidersEvent,
		payload,
		plugin_dtos.GetModelProvidersFailedEvent,
		plugin_dtos.GetModelProvidersSuccessEvent,
	)
	if err != nil {
		p.Logger.Error("failed to get list model provider", "error", err)
		return providerMap, err
	}

	if !resp.IsSuccess {
		p.Logger.Error("failed to get list model provider", "error", resp.Error)
		return providerMap, errors.New(resp.Error.Cause)
	}

	providerMap = resp.Data.List
	return providerMap, nil
}
