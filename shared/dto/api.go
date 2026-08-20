package dto

type Profile struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type ApiError struct {
	Name    string            `json:"name"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type ApiResponse[T any] struct {
	Data   T         `json:"data,omitempty"`
	Status string    `json:"status"`
	Error  *ApiError `json:"error,omitempty"`
}

type CheckHealthResponse struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Uptime int     `json:"uptime"`
	Status string  `json:"status"`
}

type Provider struct {
	ID              string   `json:"id"`
	Compatible      string   `json:"compatible"`
	Name            string   `json:"name"`
	BaseURL         string   `json:"base_url"`
	IsEnable        bool     `json:"is_enable"`
	AvailableModels []string `json:"available_models"`
	APIKey          string   `json:"api_key"`
}

type ProviderInfo struct {
	ID              string   `json:"id"`
	Compatible      string   `json:"compatible"`
	Name            string   `json:"name"`
	BaseURL         string   `json:"base_url"`
	IsEnable        bool     `json:"is_enable"`
	AvailableModels []string `json:"available_models"`
}

type GetListProviderResponse struct {
	Providers []ProviderInfo `json:"providers"`
}

type GetProviderByIDRequest struct {
	ID string `json:"id"`
}

type GetProviderByIDResponse struct {
	Provider Provider `json:"provider"`
}

type UpdateProviderRequest struct {
	ID              string   `json:"id"`
	Compatible      string   `json:"compatible"`
	Name            string   `json:"name"`
	BaseURL         string   `json:"base_url"`
	APIKey          string   `json:"api_key"`
	IsEnable        bool     `json:"is_enable"`
	AvailableModels []string `json:"available_models"`
	ModelName       string   `json:"model_name"`
}

type UpdateProviderResponse struct {
}

type CreateProviderRequest struct {
	Compatible      string   `json:"compatible"`
	Name            string   `json:"name"`
	BaseURL         string   `json:"base_url"`
	APIKey          string   `json:"api_key"`
	IsEnable        bool     `json:"is_enable"`
	AvailableModels []string `json:"available_models"`
}

type CreateProviderResponse struct {
}

type DeleteProviderRequest struct {
	ID string `json:"id"`
}

type DeleteProviderResponse struct {
}

type GeneralSettings struct {
	CurrentProvider string `json:"current_provider"`
	CurrentModel    string `json:"current_model"`
}

type GetGeneralSettingsResponse struct {
	Settings *GeneralSettings `json:"settings"`
}

type SaveGeneralSettingsRequest struct {
	Settings *GeneralSettings `json:"settings"`
}

type SaveGeneralSettingsResponse struct {
}

type FetchModelRequest struct {
	Compatible string `json:"compatible"`
	BaseURL    string `json:"base_url"`
	APIKey     string `json:"api_key"`
}

type FetchModelResponse struct {
	Models []string `json:"models"`
}

type PluginInfo struct {
	ID         string `json:"id"`
	PluginID   string `json:"plugin_id"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
	ApiVersion string `json:"api_version"`
	Author     string `json:"author"`
	Version    string `json:"version"`
}

type GetAllPluginsResponse struct {
	Plugins []PluginInfo `json:"plugins"`
}

type InstallLocalPluginRequest struct {
	Path string `json:"path"`
}

type InstallLocalPluginResponse struct {
}

type GetProfileResponse struct {
	Profile *Profile `json:"profile"`
}

type RemovePluginRequest struct {
	ID string `json:"id"`
}

type RemovePluginResponse struct {
}
