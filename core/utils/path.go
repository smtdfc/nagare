package utils

import (
	"log"
	"os"
	"path"
)

var ConfigFile = "config.json"
var LogDir = "logs"
var DataDir string

func init() {
	userDir, _ := os.UserConfigDir()
	DataDir = path.Join(userDir, ".nagare")
	ConfigFile = path.Join(DataDir, "config.json")
	LogDir = path.Join(DataDir, "logs")

	err := os.MkdirAll(DataDir, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = os.MkdirAll(LogDir, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// fmt.Println(ConfigFile)
}
