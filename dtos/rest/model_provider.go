package rest

const (
	ListModelProvidersEndpoint      = "/api/v1/user/model-providers/list"
	GetModelProviderDetailsEndpoint = "/api/v1/user/model-providers/details"
	AddModelProviderEndpoint        = "/api/v1/user/model-providers/add"
	DeleteModelProviderEndpoint     = "/api/v1/user/model-providers/delete"
	GetModelProviderModelsEndpoint  = "/api/v1/user/model-providers/models"
)

type ModelProviderInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Compatible string `json:"compatible"`
}

type ModelProvider struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Compatible string      `json:"compatible"`
	Models     []string    `json:"models"`
	BaseURL    string      `json:"baseUrl"`
	Credential *Credential `json:"credential,omitempty"`
}

type GetListModelProviderResponse struct {
	Providers []*ModelProvider `json:"providers"`
}

type GetModelProviderDetailsResponse struct {
	Provider *ModelProvider `json:"provider"`
}

type AddModelProviderRequest struct {
	Name         string   `json:"name"`
	Compatible   string   `json:"compatible"`
	ApiKey       string   `json:"apiKey"`
	Models       []string `json:"models"`
	BaseURL      string   `json:"baseUrl"`
	CredentialID string   `json:"credentialId"`
}

type AddModelProviderResponse struct {
	Provider *ModelProvider `json:"provider"`
}

type DeleteModelProviderRequest struct {
	ID string `json:"id"`
}

type GetModelProviderModelsRequest struct {
	ID string `json:"id"`
}

type GetModelProviderModelsResponse struct {
	Models []string `json:"models"`
}
