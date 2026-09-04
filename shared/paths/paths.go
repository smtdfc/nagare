package paths

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/smtdfc/nagare/shared/helpers"
)

var UserHomeDir = ""
var DataDir = ""
var ConfigFile = ""
var ConfigDir = ""
var GatewayBinFile = ""
var LogDir = ""
var PluginLogDir = ""
var PluginDir = ""
var PluginConfigDir = ""
var DatabaseDir = ""
var TempDir = ""

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("Could not determine the user's home directory: %w", err))
	}

	DataDir = filepath.Join(home, ".nagare")
	ConfigDir = filepath.Join(DataDir, "configs")
	DatabaseDir = filepath.Join(DataDir, "databases")
	LogDir = filepath.Join(DataDir, "logs")
	PluginLogDir = filepath.Join(LogDir, "plugins")
	PluginDir = filepath.Join(DataDir, "plugins")
	TempDir = filepath.Join(DataDir, "temp")
	PluginConfigDir = filepath.Join(ConfigDir, "plugin")
	ConfigFile = filepath.Join(DataDir, "config.json")
	GatewayBinFile = filepath.Join("/opt", "nagare", "nagare-gateway")
	paths := []string{
		DataDir,
		ConfigDir,
		DatabaseDir,
		LogDir,
		PluginLogDir,
		PluginDir,
		TempDir,
		PluginConfigDir,
	}

	for _, p := range paths {
		err := helpers.EnsurePathExist(p)
		if err != nil {
			println(err)
		}
	}
}
