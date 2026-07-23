package config

import (
	"encoding/json"
	"os"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/shared/paths"
)

type ConfigIO struct{}

func (c *ConfigIO) Read() (*domains.Config, error) {
	var conf domains.Config
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

func (c *ConfigIO) Write(conf *domains.Config) error {
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

func NewConfigIO() *ConfigIO {
	return &ConfigIO{}
}
