package manager

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	llm_provider "github.com/smtdfc/nagare/core/llm/provider"
	"github.com/smtdfc/nagare/core/llm/provider/adapters"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
)

type LLMProviderManager struct {
	llmProviderRepo   *repositories.LLMProviderRepository
	credentialRepo    *repositories.CredentialRepository
	llmProviderMapper *mappers.LLMProviderMapper
	logger            *logger.BaseLogger
	adapterLogger     *logger.BaseLogger
}

func (l *LLMProviderManager) GetAllProvider(ctx context.Context) ([]*llm_provider.LLMProviderConfig, error) {
	providers, err := l.llmProviderRepo.FindAll(ctx)
	if err != nil {
		return nil, custom_errors.ErrGetAllLLMProviderFailed
	}

	return l.llmProviderMapper.ToDomains(providers), nil
}

func (l *LLMProviderManager) GetProviderByID(ctx context.Context, id string) (*llm_provider.LLMProviderConfig, error) {
	provider, err := l.llmProviderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, custom_errors.ErrGetLLMProviderFailed
	}

	if provider == nil {
		return nil, custom_errors.ErrLLMProviderNotFound
	}

	return l.llmProviderMapper.ToDomain(provider), nil
}

func (l *LLMProviderManager) AddProvider(ctx context.Context, name, baseURL, compatible, apiKey, credential string, models []string) (*llm_provider.LLMProviderConfig, error) {
	conf := &llm_provider.LLMProviderConfig{
		Name:       name,
		Compatible: llm_provider.GetCompatibleFromString(compatible),
		ApiKey:     apiKey,
		Models:     models,
		BaseURL:    baseURL,
	}
	if credential != "" {
		credentialEntity, err := l.credentialRepo.FindByID(ctx, credential)
		if err != nil {
			return nil, custom_errors.ErrAddLLMProviderFailed
		}
		if credentialEntity == nil {
			return nil, custom_errors.ErrCredentialNotFound
		}
		conf.CredentialID = credentialEntity.ID
	}

	provider, err := l.llmProviderRepo.Add(ctx, l.llmProviderMapper.ToEntity(conf))
	if err != nil {
		return nil, custom_errors.ErrAddLLMProviderFailed
	}

	return l.llmProviderMapper.ToDomain(provider), nil
}

func (l *LLMProviderManager) DeleteProvider(ctx context.Context, id string) error {
	err := l.llmProviderRepo.DeleteByID(ctx, id)
	if err != nil {
		return custom_errors.ErrDeleteLLMProviderFailed
	}

	return nil
}

func (l *LLMProviderManager) FetchAvailableModels(ctx context.Context, id string) ([]string, error) {
	conf, err := l.llmProviderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, custom_errors.ErrLLMProviderNotFound
	}

	adapter, err := l.GetAdapter(l.llmProviderMapper.ToDomain(conf))
	if err != nil {
		return nil, err
	}

	if adapter == nil {
		return nil, custom_errors.ErrLLMProviderNotSupported
	}

	models, err := adapter.GetModels(ctx)
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (l *LLMProviderManager) GetAdapter(provider *llm_provider.LLMProviderConfig) (llm_provider.LLMProviderAdapter, error) {
	switch provider.Compatible {
	case llm_provider.OpenAICompatible:
		apiKey := provider.ApiKey
		if provider.ApiKey == "" && provider.Credential != nil {
			apiKey = provider.Credential.ApiKey
		}

		return adapters.NewOpenAICompatibleAdapter(
			provider.BaseURL,
			apiKey,
			provider.Models,
			l.adapterLogger.Clone(),
		), nil
	}
	return nil, custom_errors.ErrLLMProviderNotSupported
}

// @Injectable
func NewLLMProviderManager(llmProviderRepo *repositories.LLMProviderRepository, credentialRepo *repositories.CredentialRepository, llmProviderMapper *mappers.LLMProviderMapper, logger *logger.BaseLogger) *LLMProviderManager {
	return &LLMProviderManager{
		llmProviderRepo:   llmProviderRepo,
		credentialRepo:    credentialRepo,
		logger:            logger.With("module", "llm-manager"),
		llmProviderMapper: llmProviderMapper,
		adapterLogger:     logger.Clone(),
	}
}
