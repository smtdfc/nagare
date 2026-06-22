package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/utils"
)

var appLogger = logger.GetLogger()

func LoadConfig() (*Config, error) {

	data, err := os.ReadFile(utils.ConfigFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Return default config
			conf := &Config{

				Providers: map[string]Provider{},
			}
			err = SaveConfig(conf)
			return conf, err
		}
		appLogger.Error("Configuration file loading failed.", "error", err)
		return nil, err
	}

	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		appLogger.Error("Configuration file loading failed.", "error", err)
		return nil, err
	}
	if conf.Providers == nil {
		conf.Providers = make(map[string]Provider)
	}

	appLogger.Info("Configuration file loaded successfully.")
	return &conf, nil
}

func SaveConfig(conf *Config) error {
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		appLogger.Error("Saving the configuration file failed.", "error", err)
		return err
	}

	err = os.WriteFile(utils.ConfigFile, data, 0644)
	if err != nil {
		appLogger.Error("Saving the configuration file failed.", "error", err)
		return err
	}

	appLogger.Info("Configuration file saved successfully..")
	return nil
}
