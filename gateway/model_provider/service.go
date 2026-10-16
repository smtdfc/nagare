package model_provider

import (
	"context"

	model_provider_mgr "github.com/smtdfc/nagare/core/model_provider/manager"
	"github.com/smtdfc/nagare/core/model_provider/provider"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

func toModelProviderDTO(domain *provider.ModelProviderConfig) *rest.ModelProvider {
	var credentialDTO *rest.Credential
	if domain.Credential != nil {
		credentialDTO = &rest.Credential{
			ID:   domain.Credential.ID.String(),
			Name: domain.Credential.Name,
		}

	}

	return &rest.ModelProvider{
		ID:         domain.ID.String(),
		Name:       domain.Name,
		Compatible: domain.Compatible.ToString(),
		Models:     domain.Models,
		BaseURL:    domain.BaseURL,
		Credential: credentialDTO,
	}
}

type ModelProviderService struct {
	modelProviderMgr *model_provider_mgr.ModelProviderManager
}

func (l *ModelProviderService) ListProviders(ctx context.Context) (*rest.GetListModelProviderResponse, error) {
	providers, err := l.modelProviderMgr.GetAllProvider(ctx)
	if err != nil {
		return nil, err
	}

	return &rest.GetListModelProviderResponse{
		Providers: helpers.SliceMap(providers, toModelProviderDTO),
	}, nil
}

func (l *ModelProviderService) GetProviderDetails(ctx context.Context, providerID string) (*rest.GetModelProviderDetailsResponse, error) {
	modelProvider, err := l.modelProviderMgr.GetProviderByID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	return &rest.GetModelProviderDetailsResponse{
		Provider: toModelProviderDTO(modelProvider),
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
		Provider: toModelProviderDTO(modelProvider),
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
) *ModelProviderService {
	return &ModelProviderService{
		modelProviderMgr: modelProviderMgr,
	}
}
