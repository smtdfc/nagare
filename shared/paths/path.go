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
}
