package manager

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	model_provider_mapper "github.com/smtdfc/nagare/core/mappers/model_provider"
	model_provider "github.com/smtdfc/nagare/core/model_provider/provider"
	"github.com/smtdfc/nagare/core/model_provider/provider/adapters"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
)

type ModelProviderManager struct {
	modelProviderRepo   *repositories.ModelProviderRepository
	credentialRepo      *repositories.CredentialRepository
	modelProviderMapper model_provider_mapper.IModelProviderMapper
	logger              *logger.BaseLogger
	adapterLogger       *logger.BaseLogger
}

func (m *ModelProviderManager) GetAllProvider(ctx context.Context) ([]*model_provider.ModelProviderConfig, error) {
	providers, err := m.modelProviderRepo.FindAll(ctx)
	if err != nil {
		return nil, custom_errors.ErrGetAllModelProviderFailed
	}

	return m.modelProviderMapper.ToDomains(providers), nil
}

func (m *ModelProviderManager) GetProviderByID(ctx context.Context, id string) (*model_provider.ModelProviderConfig, error) {
	provider, err := m.modelProviderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, custom_errors.ErrGetModelProviderFailed
	}

	if provider == nil {
		return nil, custom_errors.ErrModelProviderNotFound
	}

	return m.modelProviderMapper.ToDomain(provider), nil
}

func (m *ModelProviderManager) AddProvider(ctx context.Context, name, baseURL, compatible, apiKey, credential string, models []string) (*model_provider.ModelProviderConfig, error) {
	conf := &model_provider.ModelProviderConfig{
		Name:       name,
		Compatible: model_provider.GetCompatibleFromString(compatible),
		ApiKey:     apiKey,
		Models:     models,
		BaseURL:    baseURL,
	}
	if credential != "" {
		credentialEntity, err := m.credentialRepo.FindByID(ctx, credential)
		if err != nil {
			return nil, custom_errors.ErrAddModelProviderFailed
		}
		if credentialEntity == nil {
			return nil, custom_errors.ErrCredentialNotFound
		}
		conf.CredentialID = credentialEntity.ID
	}

	provider, err := m.modelProviderRepo.Add(ctx, m.modelProviderMapper.ToEntity(conf))
	if err != nil {
		return nil, custom_errors.ErrAddModelProviderFailed
	}

	return m.modelProviderMapper.ToDomain(provider), nil
}

func (m *ModelProviderManager) DeleteProvider(ctx context.Context, id string) error {
	err := m.modelProviderRepo.DeleteByID(ctx, id)
	if err != nil {
		return custom_errors.ErrDeleteModelProviderFailed
	}

	return nil
}

func (m *ModelProviderManager) FetchAvailableModels(ctx context.Context, id string) ([]string, error) {
	conf, err := m.modelProviderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, custom_errors.ErrModelProviderNotFound
	}

	adapter, err := m.GetAdapter(m.modelProviderMapper.ToDomain(conf))
	if err != nil {
		return nil, err
	}

	if adapter == nil {
		return nil, custom_errors.ErrModelProviderNotSupported
	}

	models, err := adapter.GetModels(ctx)
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (m *ModelProviderManager) GetAdapter(provider *model_provider.ModelProviderConfig) (model_provider.ModelProviderAdapter, error) {

	switch provider.Compatible {
	case model_provider.OpenAICompatible:
		apiKey := provider.ApiKey
		if provider.ApiKey == "" && provider.Credential != nil {
			apiKey = provider.Credential.ApiKey
		}
		return adapters.NewOpenAICompatibleAdapter(
			provider.BaseURL,
			apiKey,
			provider.Models,
			m.adapterLogger.Clone(),
		), nil
	}
	return nil, custom_errors.ErrModelProviderNotSupported
}

// @Injectable
func NewModelProviderManager(modelProviderRepo *repositories.ModelProviderRepository, credentialRepo *repositories.CredentialRepository, modelProviderMapper model_provider_mapper.IModelProviderMapper, logger *logger.BaseLogger) *ModelProviderManager {
	return &ModelProviderManager{
		modelProviderRepo:   modelProviderRepo,
		credentialRepo:      credentialRepo,
		logger:              logger.With("module", "model-manager"),
		modelProviderMapper: modelProviderMapper,
		adapterLogger:       logger.Clone(),
	}
}
