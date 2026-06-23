package config

import "strings"

type Mode string

func (m Mode) String() string {
	return strings.Title(string(m))
}

const (
	PROVIDER_MODE Mode = "provider"
	PROXY_MODE    Mode = "proxy"
)

type ProviderCompatible string

const (
	OPEN_AI_COMPATIBLE ProviderCompatible = "openai"
)

var AllProviderCompatible = map[string]ProviderCompatible{
	"Open AI Compatible": OPEN_AI_COMPATIBLE,
}

type Provider struct {
	Name       string             `json:"name,omitempty"`
	BaseURL    string             `json:"base_url,omitempty"`
	Compatible ProviderCompatible `json:"compatible,omitempty"`
	APIKey     string             `json:"api_key,omitempty"`
	ModelName  string             `json:"model_name"`
	Enabled    bool               `json:"enabled"`
}

type Plugin struct {
	Name     string   `json:"name"`
	Id       string   `json:"id"`
	Author   string   `json:"author"`
	Version  string   `json:"version"`
	Path     string   `json:"path"`
	Features []string `json:"features"`
	Bin      string   `json:"bin"`
}

type Config struct {
	CurrentProvider string              `json:"current_provider,omitempty"`
	Providers       map[string]Provider `json:"providers,omitempty"`
	Plugins         map[string]Plugin   `json:"plugins,omitempty"`
}
