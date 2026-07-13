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

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}
