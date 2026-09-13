package rest

const (
	ListLLMProvidersEndpoint      = "/api/v1/user/llm-providers/list"
	GetLLMProviderDetailsEndpoint = "/api/v1/user/llm-providers/details"
	AddLLMProviderEndpoint        = "/api/v1/user/llm-providers/add"
	DeleteLLMProviderEndpoint     = "/api/v1/user/llm-providers/delete"
	GetLLMProviderModelsEndpoint  = "/api/v1/user/llm-providers/models"
)

type LLMProvider struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Compatible string   `json:"compatible"`
	ApiKey     string   `json:"apiKey"`
	Models     []string `json:"models"`
	BaseURL    string   `json:"baseUrl"`
}

type GetListLLMProviderResponse struct {
	Providers []*LLMProvider `json:"providers"`
}

type GetLLMProviderDetailsResponse struct {
	Provider *LLMProvider `json:"provider"`
}

type AddLLMProviderRequest struct {
	Name       string   `json:"name"`
	Compatible string   `json:"compatible"`
	ApiKey     string   `json:"apiKey"`
	Models     []string `json:"models"`
	BaseURL    string   `json:"baseUrl"`
}

type AddLLMProviderResponse struct {
	Provider *LLMProvider `json:"provider"`
}

type DeleteLLMProviderRequest struct {
	ID string `json:"id"`
}

type GetLLMProviderModelsRequest struct {
	ID string `json:"id"`
}

type GetLLMProviderModelsResponse struct {
	Models []string `json:"models"`
}
