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
var PluginLogDir = ""
var PluginDir = ""
var PluginConfigDir = ""
var DatabaseDir = ""
var TempDir = ""
var PluginSocketPath = ""
var UploadDir = ""
var CredentialsDir = ""
var ModulesDir = ""
var MemoryDir = ""
var VectorIndexFile = ""

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("could not determine the user's home directory: %w", err))
	}

	DataDir = filepath.Join(home, ".nagare")
	ModulesDir = filepath.Join(DataDir, "modules")
	MemoryDir = filepath.Join(DataDir, "memory")
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
	CredentialsDir = filepath.Join(DataDir, "credentials")
	VectorIndexFile = filepath.Join(MemoryDir, "vector_index.tvim")
	if runtime.GOOS == "windows" {
		PluginSocketPath = `\\.\pipe\nagare-plugin`
	} else {
		PluginSocketPath = filepath.Join(DataDir, "nagare.sock")
	}
	paths := []string{
		DataDir,
		ConfigDir,
		MemoryDir,
		ModulesDir,
		UploadDir,
		DatabaseDir,
		LogDir,
		PluginLogDir,
		PluginDir,
		TempDir,
		PluginConfigDir,
		CredentialsDir,
	}

	for _, p := range paths {
		err := helpers.EnsurePathExist(p)
		if err != nil {
			println(err)
		}
	}
}
