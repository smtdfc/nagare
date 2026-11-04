package sdk

import (
	"context"
	"errors"

	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func (p *PluginClient) GetListModelProvider(ctx context.Context) (map[string]string, error) {
	payload := contracts.GetModelProvidersEventPayload{}
	var providerMap = make(map[string]string)
	resp, err := sendAndWait[contracts.GetModelProvidersSuccessEventPayload, contracts.GetModelProvidersFailedEventPayload](
		p,
		ctx,
		contracts.GetModelProvidersEvent,
		payload,
		contracts.GetModelProvidersFailedEvent,
		contracts.GetModelProvidersSuccessEvent,
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
