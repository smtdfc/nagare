package config

import (
	"os"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/paths"
)

type ConfigIO struct{}

func (c *ConfigIO) Read() (*domains.Config, error) {
	raw, err := os.ReadFile(paths.ConfigFile)
	if err != nil {
		return nil, custom_errors.NewConfigError("Configuration loading error: The specified config file could not be read or is invalid.")
	}

	conf, err := helpers.FromJson[domains.Config](string(raw))
	if err != nil {
		return nil, custom_errors.NewConfigError("Configuration loading error: The specified config file could not be read or is invalid.")
	}

	return conf, nil
}

func (c *ConfigIO) Write(conf *domains.Config) error {
	raw, err := helpers.MapObjectToJson(conf)
	if err != nil {
		return custom_errors.NewConfigError("Configuration saving error: Could not write to the specified config file.")
	}

	err = os.WriteFile(paths.ConfigFile, []byte(raw), 0644)
	if err != nil {
		return custom_errors.NewConfigError("Configuration saving error: Could not write to the specified config file.")
	}

	return nil
}

func NewConfigIO() *ConfigIO {
	return &ConfigIO{}
}
