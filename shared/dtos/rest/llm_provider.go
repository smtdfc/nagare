package rest

const (
	LLMProviderGetListEndpoint            = "/api/v1/user/llm-providers/list"
	LLMProviderGetDetailsEndpoint         = "/api/v1/user/llm-providers/details"
	LLMProviderAddEndpoint                = "/api/v1/user/llm-providers/add"
	LLMProviderDeleteEndpoint             = "/api/v1/user/llm-providers/delete"
	LLMProviderGetAvailableModelsEndpoint = "/api/v1/user/llm-providers/available-models"
)

type LLMProvider struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Compatible string   `json:"compatible"`
	ApiKey     string   `json:"api_key"`
	Models     []string `json:"models"`
	BaseURL    string   `json:"base_url"`
}

type LLMProviderGetListResponse struct {
	Providers []*LLMProvider `json:"providers"`
}

type LLMProviderGetDetailsResponse struct {
	Provider *LLMProvider `json:"provider"`
}

type LLMProviderAddRequest struct {
	Name       string   `json:"name"`
	Compatible string   `json:"compatible"`
	ApiKey     string   `json:"api_key"`
	Models     []string `json:"models"`
	BaseURL    string   `json:"base_url"`
}

type LLMProviderAddResponse struct {
	Provider *LLMProvider `json:"provider"`
}

type LLMProviderDeleteRequest struct {
	ID string `json:"id"`
}

type LLMProviderGetAvailableModelsRequest struct {
	ID string `json:"id"`
}

type LLMProviderGetAvailableModelsResponse struct {
	Models []string `json:"models"`
}
