package model_provider

import (
	"context"

	model_provider_mgr "github.com/smtdfc/nagare/core/model_provider/manager"
	"github.com/smtdfc/nagare/dtos/rest"
)

type ModelProviderService struct {
	modelProviderMgr       *model_provider_mgr.ModelProviderManager
	modelProviderDtoMapper IModelProviderDtoMapper
}

func (l *ModelProviderService) ListProviders(ctx context.Context) (*rest.GetListModelProviderResponse, error) {
	providers, err := l.modelProviderMgr.GetAllProvider(ctx)
	if err != nil {
		return nil, err
	}

	return &rest.GetListModelProviderResponse{
		Providers: l.modelProviderDtoMapper.ToModelProviderDTOs(providers),
	}, nil
}

func (l *ModelProviderService) GetProviderDetails(ctx context.Context, providerID string) (*rest.GetModelProviderDetailsResponse, error) {
	modelProvider, err := l.modelProviderMgr.GetProviderByID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	return &rest.GetModelProviderDetailsResponse{
		Provider: l.modelProviderDtoMapper.ToModelProviderDTO(modelProvider),
	}, nil
}

func (l *ModelProviderService) AddProvider(ctx context.Context, request *rest.AddModelProviderRequest) (*rest.AddModelProviderResponse, error) {
	modelProvider, err := l.modelProviderMgr.AddProvider(
		ctx,
		request.Name,
		request.BaseURL,
		request.Compatible,
		request.ApiKey,
		request.CredentialID,
		request.Models,
	)

	if err != nil {
		return nil, err
	}

	return &rest.AddModelProviderResponse{
		Provider: l.modelProviderDtoMapper.ToModelProviderDTO(modelProvider),
	}, nil
}

func (l *ModelProviderService) DeleteProvider(ctx context.Context, request *rest.DeleteModelProviderRequest) error {
	return l.modelProviderMgr.DeleteProvider(ctx, request.ID)
}

func (l *ModelProviderService) GetModels(ctx context.Context, request *rest.GetModelProviderModelsRequest) (*rest.GetModelProviderModelsResponse, error) {
	models, err := l.modelProviderMgr.FetchAvailableModels(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &rest.GetModelProviderModelsResponse{
		Models: models,
	}, nil
}

// @Injectable
func NewModelProviderService(
	modelProviderMgr *model_provider_mgr.ModelProviderManager,
	modelProviderDtoMapper IModelProviderDtoMapper,
) *ModelProviderService {
	return &ModelProviderService{
		modelProviderMgr:       modelProviderMgr,
		modelProviderDtoMapper: modelProviderDtoMapper,
	}
}
