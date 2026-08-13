package config

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/persistence/database/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
)

type ConfigManager struct {
	llmProviderRepo *repositories.LLMProviderRepository
	kvRepo          *repositories.KVRepository
}

func (c *ConfigManager) GetGeneralConfig() (*domains.GeneralConfig, error) {
	mapper := &mappers.KVMapper{}
	kv, err := c.kvRepo.GetAllKeyByTarget("nagare_general_config")
	if err != nil {
		return nil, err
	}

	return mapper.ToGeneralConfig(kv), nil
}

func (c *ConfigManager) SaveGeneralConfig(config *domains.GeneralConfig) error {
	mapper := &mappers.KVMapper{}
	kv := mapper.FromGeneralConfig(config, "nagare_general_config")
	return c.kvRepo.Save(kv)
}

func (c *ConfigManager) GetLLMProviderConfigByID(id string) (*domains.LLMProviderConfig, error) {
	mapper := &mappers.LLMProviderMapper{}
	llmProvider, err := c.llmProviderRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToConfig(llmProvider), nil
}

func (c *ConfigManager) SaveLLMProviderConfig(config *domains.LLMProviderConfig) error {
	mapper := &mappers.LLMProviderMapper{}
	llmProvider := mapper.ToModel(config)
	return c.llmProviderRepo.Save(llmProvider)
}

func (c *ConfigManager) GetListProviders() ([]*domains.LLMProviderConfigInfo, error) {
	mapper := &mappers.LLMProviderMapper{}
	providers, err := c.llmProviderRepo.FindAll()
	if err != nil {
		return nil, err
	}

	infos := make([]*domains.LLMProviderConfigInfo, 0, len(providers))
	for _, provider := range providers {
		infos = append(infos, mapper.ToInfo(&provider))
	}
	return infos, nil
}

func NewConfigManager(llmProviderRepo *repositories.LLMProviderRepository, kvRepo *repositories.KVRepository) *ConfigManager {
	return &ConfigManager{
		llmProviderRepo: llmProviderRepo,
		kvRepo:          kvRepo,
	}
}
