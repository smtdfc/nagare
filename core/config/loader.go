package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
)

var ConfigFile = "config.json"
var DataDir string

func init() {
	userDir, _ := os.UserConfigDir()
	DataDir = path.Join(userDir, ".nagare")
	ConfigFile = path.Join(DataDir, "config.json")

	err := os.MkdirAll(DataDir, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(ConfigFile)
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Return default config
			conf := &Config{
				CurrentMode: PROVIDER_MODE,
				Providers:   map[string]Provider{},
			}
			err = SaveConfig(conf)
			return conf, err
		}
		return nil, err
	}

	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	if conf.Providers == nil {
		conf.Providers = make(map[string]Provider)
	}
	return &conf, nil
}

func SaveConfig(conf *Config) error {
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, data, 0644)
}
