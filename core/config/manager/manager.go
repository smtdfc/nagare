package manager

import (
	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/persistence/database/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
)

type ConfigManager struct {
	llmProviderRepo *repositories.LLMProviderRepository
	kvRepo          *repositories.KVRepository
}

func (c *ConfigManager) GetGeneralConfig() (*config.GeneralConfig, error) {
	mapper := &mappers.KVMapper{}
	kv, err := c.kvRepo.GetAllKeyByTarget("nagare_general_config")
	if err != nil {
		ConfigManagerLogger.Error("Failed to get general config", "error", err)
		return nil, custom_errors.NewConfigError("Failed to get general config")
	}

	return mapper.ToGeneralConfig(kv), nil
}

func (c *ConfigManager) SaveGeneralConfig(config *config.GeneralConfig) error {
	mapper := &mappers.KVMapper{}
	kv := mapper.FromGeneralConfig(config, "nagare_general_config")
	if err := c.kvRepo.Save(kv); err != nil {
		ConfigManagerLogger.Error("Failed to save general config", "error", err)
		return custom_errors.NewConfigError("Failed to save general config")
	}
	return nil
}

func (c *ConfigManager) GetLLMProviderConfigByID(id string) (*config.LLMProviderConfig, error) {
	mapper := &mappers.LLMProviderMapper{}
	llmProvider, err := c.llmProviderRepo.FindByID(id)
	if err != nil {
		ConfigManagerLogger.Error("Failed to get LLM provider config", "error", err)
		return nil, custom_errors.NewConfigError("Failed to get LLM provider config")
	}
	return mapper.ToConfig(llmProvider), nil
}

func (c *ConfigManager) SaveLLMProviderConfig(config *config.LLMProviderConfig) error {
	mapper := &mappers.LLMProviderMapper{}
	llmProvider := mapper.ToModel(config)
	if err := c.llmProviderRepo.UpdateByID(config.ID, llmProvider); err != nil {
		ConfigManagerLogger.Error("Failed to save LLM provider config", "error", err)
		return custom_errors.NewConfigError("Failed to save LLM provider config")
	}
	return nil
}

func (c *ConfigManager) GetListProviders() ([]*config.LLMProviderConfigInfo, error) {
	mapper := &mappers.LLMProviderMapper{}
	providers, err := c.llmProviderRepo.FindAll()
	if err != nil {
		ConfigManagerLogger.Error("Failed to get list of LLM providers", "error", err)
		return nil, custom_errors.NewConfigError("Failed to get list of LLM providers")
	}

	infos := make([]*config.LLMProviderConfigInfo, 0, len(providers))
	for _, provider := range providers {
		infos = append(infos, mapper.ToInfo(&provider))
	}
	return infos, nil
}

func (c *ConfigManager) CreateLLMProviderConfig(config *config.LLMProviderConfig) error {
	mapper := &mappers.LLMProviderMapper{}
	llmProvider := mapper.ToModel(config)
	if err := c.llmProviderRepo.CreateProvider(llmProvider); err != nil {
		ConfigManagerLogger.Error("Failed to create LLM provider config", "error", err)
		return custom_errors.NewConfigError("Failed to create LLM provider config")
	}
	return nil
}

func (c *ConfigManager) DeleteLLMProviderConfig(id string) error {
	if err := c.llmProviderRepo.DeleteByID(id); err != nil {
		ConfigManagerLogger.Error("Failed to delete LLM provider config", "error", err)
		return custom_errors.NewConfigError("Failed to delete LLM provider config")
	}
	return nil
}

// @Injectable
func NewConfigManager(llmProviderRepo *repositories.LLMProviderRepository, kvRepo *repositories.KVRepository) *ConfigManager {
	return &ConfigManager{
		llmProviderRepo: llmProviderRepo,
		kvRepo:          kvRepo,
	}
}
