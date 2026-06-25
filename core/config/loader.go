package config

import (
	"encoding/json"
	"errors"
	"os"

	nagare_logger "github.com/smtdfc/nagare/shared/logger"
	nagare_path "github.com/smtdfc/nagare/shared/path"
)

func LoadConfig() (*Config, error) {
	var appLogger = nagare_logger.GetLogger("Config loader")

	data, err := os.ReadFile(nagare_path.ConfigFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Return default config
			conf := &Config{
				Plugins:   map[string]Plugin{},
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
	var appLogger = nagare_logger.GetLogger("Config loader")
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		appLogger.Error("Saving the configuration file failed.", "error", err)
		return err
	}

	err = os.WriteFile(nagare_path.ConfigFile, data, 0644)
	if err != nil {
		appLogger.Error("Saving the configuration file failed.", "error", err)
		return err
	}

	appLogger.Info("Configuration file saved successfully..")
	return nil
}
