package dto

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
	APIKey          string   `json:"api_key"`
	IsEnable        bool     `json:"is_enable"`
	AvailableModels []string `json:"available_models"`
}

type GetListProviderResponse struct {
	Providers []Provider `json:"providers"`
}
