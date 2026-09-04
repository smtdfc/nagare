package services

import (
	"context"

	"github.com/smtdfc/nagare/core/llm_provider"
	llm_provider_mgr "github.com/smtdfc/nagare/core/llm_provider/manager"
	"github.com/smtdfc/nagare/shared/dtos/rest"
	"github.com/smtdfc/nagare/shared/helpers"
)

func toLLMProviderDTO(domain *llm_provider.LLMProviderConfig) *rest.LLMProvider {
	return &rest.LLMProvider{
		ID:         domain.ID,
		Name:       domain.Name,
		Compatible: domain.Compatible.ToString(),
		ApiKey:     domain.ApiKey,
		Models:     domain.Models,
		BaseURL:    domain.BaseURL,
	}
}

type LLMProviderService struct {
	llmProviderMgr *llm_provider_mgr.LLMProviderManager
}

func (l *LLMProviderService) GetListProvider(ctx context.Context) (*rest.LLMProviderGetListResponse, error) {
	providers, err := l.llmProviderMgr.GetAllProvider(ctx)
	if err != nil {
		return nil, err
	}

	return &rest.LLMProviderGetListResponse{
		Providers: helpers.SliceMap(providers, toLLMProviderDTO),
	}, nil
}

func (l *LLMProviderService) GetProviderDetails(ctx context.Context, providerID string) (*rest.LLMProviderGetDetailsResponse, error) {
	provider, err := l.llmProviderMgr.GetProviderByID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	return &rest.LLMProviderGetDetailsResponse{
		Provider: toLLMProviderDTO(provider),
	}, nil
}

func (l *LLMProviderService) AddProvider(ctx context.Context, request *rest.LLMProviderAddRequest) (*rest.LLMProviderAddResponse, error) {
	provider, err := l.llmProviderMgr.AddProvider(
		ctx,
		request.Name,
		request.BaseURL,
		request.Compatible,
		request.ApiKey,
		request.Models,
	)

	if err != nil {
		return nil, err
	}

	return &rest.LLMProviderAddResponse{
		Provider: toLLMProviderDTO(provider),
	}, nil
}

func (l *LLMProviderService) DeleteProvider(ctx context.Context, request *rest.LLMProviderDeleteRequest) error {
	return l.llmProviderMgr.DeleteProvider(ctx, request.ID)
}

func (l *LLMProviderService) GetAvailableModels(ctx context.Context, request *rest.LLMProviderGetAvailableModelsRequest) (*rest.LLMProviderGetAvailableModelsResponse, error) {
	models, err := l.llmProviderMgr.FetchAvailableModels(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &rest.LLMProviderGetAvailableModelsResponse{
		Models: models,
	}, nil
}

// @Injectable
func NewLLMProviderService(
	llmProviderMgr *llm_provider_mgr.LLMProviderManager,
) *LLMProviderService {
	return &LLMProviderService{
		llmProviderMgr: llmProviderMgr,
	}
}
