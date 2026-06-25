package nagare_path

import (
	"log"
	"os"
	"path"
)

var ConfigFile = "config.json"
var LogDir = "logs"
var PluginDir = "plugins"
var TempDir = "temp"
var DataDir string

func InitGlobalPath() {
	userDir, _ := os.UserConfigDir()
	DataDir = path.Join(userDir, ".nagare")
	ConfigFile = path.Join(DataDir, "config.json")
	LogDir = path.Join(DataDir, "logs")
	PluginDir = path.Join(DataDir, "plugins")
	TempDir = path.Join(DataDir, "temp")

	err := os.MkdirAll(DataDir, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = os.MkdirAll(LogDir, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = os.MkdirAll(PluginDir, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = os.MkdirAll(TempDir, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// fmt.Println(ConfigFile)
}
