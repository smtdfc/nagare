package config

type ProviderConfig struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	IsEnable bool   `json:"is_enable"`
	Model    string `json:"model"`
}

type Config struct {
	CurrentProvider string                     `json:"current_provider"`
	Providers       map[string]*ProviderConfig `json:"providers"`
}
