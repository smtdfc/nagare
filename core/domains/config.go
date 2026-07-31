package domains

type ProviderCompatible string

const (
	OPEN_AI ProviderCompatible = "OpenAI"
)

type ProviderConfig struct {
	ID              string             `json:"id"`
	Compatible      ProviderCompatible `json:"compatible"`
	Name            string             `json:"name"`
	BaseURL         string             `json:"base_url"`
	APIKey          string             `json:"api_key"`
	IsEnable        bool               `json:"is_enable"`
	Model           string             `json:"model"`
	AvailableModels []string           `json:"available_models"`
}

type Config struct {
	CurrentProvider string                     `json:"current_provider"`
	CurrentModel    string                     `json:"current_model"`
	Providers       map[string]*ProviderConfig `json:"providers"`
}
