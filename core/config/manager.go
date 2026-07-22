package config

import (
	"encoding/json"
	"os"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/shared/paths"
)

type ConfigManager struct{}

func (c *ConfigManager) Load() (*Config, error) {
	var conf Config
	raw, err := os.ReadFile(paths.ConfigFile)
	if err != nil {
		return nil, custom_errors.NewConfigError("Configuration loading error: The specified config file could not be read or is invalid.")
	}

	err = json.Unmarshal(raw, &conf)
	if err != nil {
		return nil, custom_errors.NewConfigError("Configuration loading error: The specified config file could not be read or is invalid.")
	}

	return &conf, nil
}

func (c *ConfigManager) Save(conf *Config) error {
	raw, err := json.Marshal(conf)
	if err != nil {
		return custom_errors.NewConfigError("Configuration saving error: Could not write to the specified config file.")
	}

	err = os.WriteFile(paths.ConfigFile, raw, 0644)
	if err != nil {
		return custom_errors.NewConfigError("Configuration saving error: Could not write to the specified config file.")
	}

	return nil
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}
