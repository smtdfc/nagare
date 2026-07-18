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
		return nil, custom_errors.NewConfigError("failed to load config ")
	}

	err = json.Unmarshal(raw, &conf)
	if err != nil {
		return nil, custom_errors.NewConfigError("failed to load config ")
	}

	return &conf, nil
}

func (c *ConfigManager) Save(conf *Config) error {
	raw, err := json.Marshal(conf)
	if err != nil {
		return custom_errors.NewConfigError("failed to save config ")
	}

	err = os.WriteFile(paths.ConfigFile, raw, 0644)
	if err != nil {
		return custom_errors.NewConfigError("failed to save config ")
	}

	return nil
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}
