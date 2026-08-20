package global

import (
	llm_manager "github.com/smtdfc/nagare/core/llm/manager"
)

var GlobalLLMManager *llm_manager.LLMManager

// func GetLLMProvider() (domains.LLMProviderAdapter, error) {
// 	generalConfig, err := GlobalConfigMgr.GetGeneralConfig()
// 	if err != nil {
// 		return nil, err
// 	}

// 	currentProviderConfig, err := GlobalConfigMgr.GetLLMProviderConfigByID(generalConfig.CurrentModel)
// 	if err != nil {
// 		return nil, err
// 	}

// 	switch currentProviderConfig.Compatible {
// 	case "OpenAI":
// 		return providers.NewOpenAICompatibleProviderAdapter(
// 			currentProviderConfig.BaseURL,
// 			currentProviderConfig.APIKey,
// 			currentProviderConfig.AvailableModels,
// 		), nil
// 	}

// 	return nil, custom_errors.NewLLMProviderError("The selected provider is not compatible with Nagare")
// }
