package config

type ProviderCompatible string

const (
	OPEN_AI ProviderCompatible = "OpenAI"
)

type ProviderConfig struct {
	ID         string             `json:"id"`
	Compatible ProviderCompatible `json:"compatible"`
	Name       string             `json:"name"`
	BaseURL    string             `json:"base_url"`
	APIKey     string             `json:"api_key"`
	IsEnable   bool               `json:"is_enable"`
	Model      string             `json:"model"`
}

type Config struct {
	CurrentProvider string                     `json:"current_provider"`
	Providers       map[string]*ProviderConfig `json:"providers"`
}
