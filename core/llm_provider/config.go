package llm_provider

import "github.com/google/uuid"

type LLMProviderCompatible string

const (
	OpenAICompatible  LLMProviderCompatible = "OpenAI"
	UnknownCompatible LLMProviderCompatible = "Unknown"
)

func (l LLMProviderCompatible) ToString() string {
	return string(l)
}

func GetCompatibleFromString(s string) LLMProviderCompatible {
	switch s {
	case string(OpenAICompatible):
		return OpenAICompatible
	default:
		return UnknownCompatible
	}
}

type LLMProviderConfig struct {
	ID         uuid.UUID
	Name       string
	Compatible LLMProviderCompatible
	ApiKey     string
	Models     []string
	BaseURL    string
}
