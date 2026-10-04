package paths

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/smtdfc/nagare/pkgs/helpers"
)

var UserHomeDir = ""
var DataDir = ""
var ConfigFile = ""
var GatewayPIDFile = ""
var ConfigDir = ""
var GatewayBinFile = ""
var LogDir = ""
var PluginLogDir = ""
var PluginDir = ""
var PluginConfigDir = ""
var DatabaseDir = ""
var TempDir = ""
var PluginSocketPath = ""
var UploadDir = ""

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("could not determine the user's home directory: %w", err))
	}

	DataDir = filepath.Join(home, ".nagare")
	UploadDir = filepath.Join(DataDir, "uploads")
	ConfigDir = filepath.Join(DataDir, "configs")
	DatabaseDir = filepath.Join(DataDir, "databases")
	LogDir = filepath.Join(DataDir, "logs")
	PluginLogDir = filepath.Join(LogDir, "plugins")
	PluginDir = filepath.Join(DataDir, "plugins")
	TempDir = filepath.Join(DataDir, "temp")
	PluginConfigDir = filepath.Join(ConfigDir, "plugins")
	ConfigFile = filepath.Join(DataDir, "config.json")
	GatewayPIDFile = filepath.Join(DataDir, ".gateway.pid")
	GatewayBinFile = filepath.Join("/opt", "nagare", "nagare-gateway")
	if runtime.GOOS == "windows" {
		PluginSocketPath = `\\.\pipe\nagare-plugin`
	} else {
		PluginSocketPath = filepath.Join(DataDir, "nagare.sock")
	}
	paths := []string{
		DataDir,
		ConfigDir,
		UploadDir,
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
