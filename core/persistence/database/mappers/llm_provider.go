package mappers

import (
	"strings"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/persistence/database/models"
)

type LLMProviderMapper struct {
}

func (m *LLMProviderMapper) ToConfig(model *models.LLMProvider) *domains.LLMProviderConfig {
	return &domains.LLMProviderConfig{
		ID:              model.ID.String(),
		Compatible:      model.Compatible,
		Name:            model.Name,
		BaseURL:         model.BaseURL,
		APIKey:          model.APIKey,
		IsEnable:        model.IsEnable,
		ModelName:       model.ModelName,
		AvailableModels: strings.Split(model.AvailableModels, ","),
	}
}

func (m *LLMProviderMapper) ToInfo(model *models.LLMProvider) *domains.LLMProviderConfigInfo {
	return &domains.LLMProviderConfigInfo{
		ID:              model.ID.String(),
		Compatible:      model.Compatible,
		Name:            model.Name,
		BaseURL:         model.BaseURL,
		IsEnable:        model.IsEnable,
		AvailableModels: strings.Split(model.AvailableModels, ","),
	}
}

func (m *LLMProviderMapper) ToModel(config *domains.LLMProviderConfig) *models.LLMProvider {
	return &models.LLMProvider{
		Compatible:      config.Compatible,
		Name:            config.Name,
		BaseURL:         config.BaseURL,
		APIKey:          config.APIKey,
		IsEnable:        config.IsEnable,
		ModelName:       config.ModelName,
		AvailableModels: strings.Join(config.AvailableModels, ","),
	}
}
