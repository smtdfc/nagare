package metadata

import (
	"errors"
	"strings"
)

type PluginBinaryMetadata map[string]string
type PluginFeatures []string

type PluginMetadata struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Version    string               `json:"version"`
	ApiVersion string               `json:"api_version"`
	Author     string               `json:"author"`
	Bin        PluginBinaryMetadata `json:"bin"`
	Features   PluginFeatures       `json:"features"`
}

func (m *PluginMetadata) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return errors.New("plugin metadata 'id' is required")
	}
	if strings.TrimSpace(m.Name) == "" {
		return errors.New("plugin metadata 'name' is required")
	}
	if strings.TrimSpace(m.Version) == "" {
		return errors.New("plugin metadata 'version' is required")
	}
	if strings.TrimSpace(m.ApiVersion) == "" {
		return errors.New("plugin metadata 'api_version' is required")
	}
	if strings.TrimSpace(m.Author) == "" {
		return errors.New("plugin metadata 'author' is required")
	}
	return nil
}
