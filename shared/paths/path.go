package paths

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var UserHomeDir = ""
var DataDir = ""
var ConfigFile = ""
var LogDir = ""
var PluginLogDir = ""
var DatabaseDir = ""

func EnsureConfigWithDefaults() error {
	if _, err := os.Stat(DataDir); os.IsNotExist(err) {
		err := os.MkdirAll(DataDir, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		defaultConfig := []byte(`{}`)

		err := os.WriteFile(ConfigFile, defaultConfig, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get home dir: %w", err))
	}

	DataDir = filepath.Join(home, ".nagare")

	dirs := []string{
		DataDir,
		filepath.Join(DataDir, "databases"),
		filepath.Join(DataDir, "logs"),
		filepath.Join(DataDir, "logs", "plugins"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(fmt.Errorf("failed to create dir %s: %w", dir, err))
		}
	}

	ConfigFile = filepath.Join(DataDir, "config.json")
	err = EnsureConfigWithDefaults()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create config file %s: %w", ConfigFile, err))
	}
}
