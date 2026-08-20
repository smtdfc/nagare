package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/smtdfc/nagare/shared/paths"
)

func ReadPluginConfig[T any](client *PluginClient) (*T, error) {
	var conf T

	pluginConfigPath := filepath.Join(paths.PluginConfigDir, fmt.Sprintf("%s.json", client.name))
	raw, err := os.ReadFile(pluginConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if writeErr := WritePluginConfig(client, &conf); writeErr != nil {
				return nil, fmt.Errorf("failed to create default plugin config file: %w", writeErr)
			}
			return &conf, nil
		}

		return nil, fmt.Errorf("failed to read plugin config file: %w", err)
	}

	err = json.Unmarshal(raw, &conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin config: %w", err)
	}

	return &conf, nil
}

func WritePluginConfig[T any](client *PluginClient, conf *T) error {
	pluginConfigPath := filepath.Join(paths.PluginConfigDir, fmt.Sprintf("%s.json", client.name))

	// dir := filepath.Dir(pluginConfigPath)
	// if err := os.MkdirAll(dir, 0755); err != nil {
	// 	return fmt.Errorf("failed to create config directory: %w", err)
	// }

	raw, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal plugin config: %w", err)
	}

	err = os.WriteFile(pluginConfigPath, raw, 0644)
	if err != nil {
		return fmt.Errorf("failed to write plugin config file: %w", err)
	}

	return nil
}
