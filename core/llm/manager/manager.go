package llm_manager

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/llm/providers"
)

type LLMManger struct {
}

func (l *LLMManger) getOpenAIModels(baseURL string, apiKey string) ([]string, error) {
	openAIAdapter := providers.NewOpenAICompatibleProviderAdapter(baseURL, apiKey, []string{})
	ctx := context.Background()
	return openAIAdapter.ListModel(ctx)
}

func (l *LLMManger) GetAvailableModels(compatable string, baseURL string, apiKey string) ([]string, error) {
	switch compatable {
	case "OpenAI":
		return l.getOpenAIModels(baseURL, apiKey)
	default:
		return nil, custom_errors.NewLLMProviderError("Not supported")
	}
}

func NewLLMManger() *LLMManger {
	return &LLMManger{}
}
