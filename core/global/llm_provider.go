package global

import (
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/domains"
	llm_manager "github.com/smtdfc/nagare/core/llm/manager"
	"github.com/smtdfc/nagare/core/llm/providers"
)

var GlobalLLMManager *llm_manager.LLMManger

func GetLLMProvider() (domains.LLMProviderAdapter, error) {
	generalConfig, err := GlobalConfigMgr.GetGeneralConfig()
	if err != nil {
		return nil, err
	}

	currentProviderConfig, err := GlobalConfigMgr.GetLLMProviderConfigByID(generalConfig.CurrentModel)
	if err != nil {
		return nil, err
	}

	switch currentProviderConfig.Compatible {
	case "OpenAI":
		return providers.NewOpenAICompatibleProviderAdapter(
			currentProviderConfig.BaseURL,
			currentProviderConfig.APIKey,
			currentProviderConfig.AvailableModels,
		), nil
	}

	return nil, custom_errors.NewLLMProviderError("The selected provider is not compatible with Nagare")
}
