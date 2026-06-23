package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/smtdfc/nagare/plugin-builder/utils"
	plugin_sdk "github.com/smtdfc/nagare/plugin-sdk"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	pluginMetadataFile := path.Join(cwd, "metadata.json")
	isExist, err := utils.FileExists(pluginMetadataFile)
	if !isExist {
		fmt.Println("File metadata.json not found")
		return
	}

	var metadata plugin_sdk.PluginMetadata
	err = utils.ReadJSON(pluginMetadataFile, &metadata)

	fmt.Printf("Building plugin %s %s\n", metadata.ID, metadata.Version)
	cmd := exec.Command("go", "build", "-v", "-o", metadata.Bin)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Build failed with error:", err)
		return
	}

	fmt.Println("Build completed successfully!")

	outputZip := fmt.Sprintf("%s.nagare_plugin", metadata.ID)
	err = createPluginPackage(metadata.Bin, pluginMetadataFile, outputZip)
	if err != nil {
		fmt.Println("Failed to create package:", err)
		return
	}

	fmt.Printf("Plugin packaged successfully: %s\n", outputZip)
}
