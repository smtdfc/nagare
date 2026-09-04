package manager

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/llm_provider"
	"github.com/smtdfc/nagare/core/llm_provider/adapters"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
)

type LLMProviderManager struct {
	llmProviderRepo   *repositories.LLMProviderRepository
	llmProviderMapper *mappers.LLMProviderMapper
	logger            *logger.BaseLogger
	adapterLogger     *logger.BaseLogger
}

func (l *LLMProviderManager) GetAllProvider(ctx context.Context) ([]*llm_provider.LLMProviderConfig, error) {
	providers, err := l.llmProviderRepo.GetAllProvider(ctx)
	if err != nil {
		return nil, custom_errors.ErrGetAllLLMProviderFailed
	}

	return l.llmProviderMapper.ToDomains(providers), nil
}

func (l *LLMProviderManager) GetProviderByID(ctx context.Context, id string) (*llm_provider.LLMProviderConfig, error) {
	provider, err := l.llmProviderRepo.GetProviderByID(ctx, id)
	if err != nil {
		return nil, custom_errors.ErrGetLLMProviderFailed
	}

	if provider == nil {
		return nil, custom_errors.ErrLLMProviderNotFound
	}

	return l.llmProviderMapper.ToDomain(provider), nil
}

func (l *LLMProviderManager) AddProvider(ctx context.Context, name, baseURL, compatible, apiKeys string, models []string) (*llm_provider.LLMProviderConfig, error) {
	conf := &llm_provider.LLMProviderConfig{
		Name:       name,
		Compatible: llm_provider.GetCompatibleFromString(compatible),
		ApiKey:     apiKeys,
		Models:     models,
		BaseURL:    baseURL,
	}

	provider, err := l.llmProviderRepo.AddProvider(ctx, l.llmProviderMapper.ToEntity(conf))
	if err != nil {
		return nil, custom_errors.ErrAddLLMProviderFailed
	}

	return l.llmProviderMapper.ToDomain(provider), nil
}

func (l *LLMProviderManager) DeleteProvider(ctx context.Context, id string) error {
	err := l.llmProviderRepo.DeleteProviderByID(ctx, id)
	if err != nil {
		return custom_errors.ErrDeleteLLMProviderFailed
	}

	return nil
}

func (l *LLMProviderManager) FetchAvailableModels(ctx context.Context, id string) ([]string, error) {
	conf, err := l.llmProviderRepo.GetProviderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, custom_errors.ErrLLMProviderNotFound
	}

	adapter, err := l.GetAdapter(l.llmProviderMapper.ToDomain(conf))
	models, err := adapter.GetModels(ctx)
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (l *LLMProviderManager) GetAdapter(provider *llm_provider.LLMProviderConfig) (llm_provider.LLMProviderAdapter, error) {
	switch provider.Compatible {
	case llm_provider.OpenAICompatible:
		return adapters.NewOpenAICompatibleAdapter(
			provider.BaseURL,
			provider.ApiKey,
			provider.Models,
			l.adapterLogger.Clone(),
		), nil
	}
	return nil, custom_errors.ErrProviderNotSupported
}

// @Injectable
func NewLLMProviderManager(llmProviderRepo *repositories.LLMProviderRepository, llmProviderMapper *mappers.LLMProviderMapper, logger *logger.BaseLogger) *LLMProviderManager {
	return &LLMProviderManager{
		llmProviderRepo:   llmProviderRepo,
		logger:            logger.With("module", "llm-manager"),
		llmProviderMapper: llmProviderMapper,
		adapterLogger:     logger.Clone(),
	}
}
