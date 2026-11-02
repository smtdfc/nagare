package paths

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/smtdfc/nagare/pkgs/helpers"
)

var DataDir = ""
var ConfigFile = ""
var GatewayPIDFile = ""
var ConfigDir = ""
var GatewayBinFile = ""
var LogDir = ""

var DatabaseDir = ""
var TempDir = ""

var UploadDir = ""
var CredentialsDir = ""

var TrustedDir = ""

var PluginDir = ""
var PluginLogDir = ""
var PluginConfigDir = ""
var PluginDataDir = ""
var PluginSocketPath = ""

var ModuleDir = ""
var ModuleLogDir = ""
var ModuleDataDir = ""
var ModuleConfigDir = ""
var ModuleSocketPath = ""

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
	TempDir = filepath.Join(DataDir, "temp")
	ConfigFile = filepath.Join(DataDir, "config.json")
	GatewayPIDFile = filepath.Join(DataDir, ".gateway.pid")
	GatewayBinFile = filepath.Join("/opt", "nagare", "nagare-gateway")
	CredentialsDir = filepath.Join(DataDir, "credentials")
	TrustedDir = filepath.Join(DataDir, "trusted")

	PluginDir = filepath.Join(DataDir, "plugins")
	PluginLogDir = filepath.Join(LogDir, "plugins")
	PluginConfigDir = filepath.Join(ConfigDir, "plugins")
	PluginDataDir = filepath.Join(DataDir, "data", "plugins")

	ModuleDir = filepath.Join(DataDir, "modules")
	ModuleLogDir = filepath.Join(LogDir, "modules")
	ModuleConfigDir = filepath.Join(ConfigDir, "modules")
	ModuleDataDir = filepath.Join(DataDir, "data", "modules")

	if runtime.GOOS == "windows" {
		PluginSocketPath = `\\.\pipe\nagare-plugin`
		ModuleSocketPath = `\\.\pipe\nagare-external-module`
	} else {
		PluginSocketPath = filepath.Join(DataDir, "nagare-plugin.sock")
		ModuleSocketPath = filepath.Join(DataDir, "nagare-external-module.sock")
	}

	paths := []string{
		DataDir,
		ConfigDir,
		UploadDir,
		DatabaseDir,
		LogDir,
		TempDir,
		CredentialsDir,
		TrustedDir,

		PluginDir,
		PluginLogDir,
		PluginConfigDir,
		PluginDataDir,

		ModuleDir,
		ModuleLogDir,
		ModuleConfigDir,
		ModuleDataDir,
	}

	for _, p := range paths {
		err := helpers.EnsurePathExist(p)
		if err != nil {
			println(err)
		}
	}
}
