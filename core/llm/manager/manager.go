package llm_manager

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/llm"
	"github.com/smtdfc/nagare/core/llm/providers"
)

type LLMManager struct {
}

func (l *LLMManager) getOpenAIModels(baseURL string, apiKey string) ([]string, error) {
	openAIAdapter := providers.NewOpenAICompatibleProviderAdapter(baseURL, apiKey, []string{})
	ctx := context.Background()
	models, err := openAIAdapter.ListModel(ctx)
	if err != nil {
		llm.LLMLogger.Error("failed to get OpenAI models", "error", err)
		return nil, custom_errors.NewLLMProviderError("failed to get models: " + err.Error())
	}
	return models, nil
}

func (l *LLMManager) GetAvailableModels(compatable string, baseURL string, apiKey string) ([]string, error) {
	switch compatable {
	case "OpenAI":
		return l.getOpenAIModels(baseURL, apiKey)
	default:
		return nil, custom_errors.NewLLMProviderError("Not supported")
	}
}

// @Injectable
func NewLLMManager() *LLMManager {
	return &LLMManager{}
}
