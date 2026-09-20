package llm_provider

import (
	"context"

	"github.com/smtdfc/nagare/core/llm_provider"
	llm_provider_mgr "github.com/smtdfc/nagare/core/llm_provider/manager"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/shared/helpers"
)

func toLLMProviderDTO(domain *llm_provider.LLMProviderConfig) *rest.LLMProvider {
	return &rest.LLMProvider{
		ID:         domain.ID.String(),
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

func (l *LLMProviderService) ListProviders(ctx context.Context) (*rest.GetListLLMProviderResponse, error) {
	providers, err := l.llmProviderMgr.GetAllProvider(ctx)
	if err != nil {
		return nil, err
	}

	return &rest.GetListLLMProviderResponse{
		Providers: helpers.SliceMap(providers, toLLMProviderDTO),
	}, nil
}

func (l *LLMProviderService) GetProviderDetails(ctx context.Context, providerID string) (*rest.GetLLMProviderDetailsResponse, error) {
	provider, err := l.llmProviderMgr.GetProviderByID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	return &rest.GetLLMProviderDetailsResponse{
		Provider: toLLMProviderDTO(provider),
	}, nil
}

func (l *LLMProviderService) AddProvider(ctx context.Context, request *rest.AddLLMProviderRequest) (*rest.AddLLMProviderResponse, error) {
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

	return &rest.AddLLMProviderResponse{
		Provider: toLLMProviderDTO(provider),
	}, nil
}

func (l *LLMProviderService) DeleteProvider(ctx context.Context, request *rest.DeleteLLMProviderRequest) error {
	return l.llmProviderMgr.DeleteProvider(ctx, request.ID)
}

func (l *LLMProviderService) GetModels(ctx context.Context, request *rest.GetLLMProviderModelsRequest) (*rest.GetLLMProviderModelsResponse, error) {
	models, err := l.llmProviderMgr.FetchAvailableModels(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &rest.GetLLMProviderModelsResponse{
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
