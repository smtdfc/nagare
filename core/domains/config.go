package domains

type GeneralConfig struct {
	CurrentModel    string `json:"current_model"`
	CurrentProvider string `json:"current_provider"`
}

type LLMProviderConfig struct {
	ID              string
	Compatible      string
	Name            string
	BaseURL         string
	APIKey          string
	IsEnable        bool
	ModelName       string
	AvailableModels []string
}

type LLMProviderConfigInfo struct {
	ID              string
	Compatible      string
	Name            string
	BaseURL         string
	IsEnable        bool
	AvailableModels []string
}
