package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/smtdfc/nagare/plugin-builder/utils"
	"github.com/smtdfc/nagare/plugin-sdk/plugin"
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

	var metadata plugin.PluginMetadata
	err = utils.ReadJSON(pluginMetadataFile, &metadata)

	fmt.Printf("Building plugin %s %s\n", metadata.ID, metadata.Version)

	for platform, output := range metadata.Bin {
		parts := strings.Split(platform, "/")
		if len(parts) != 2 {
			fmt.Println("Invalid platform format, expected os/arch:", platform)
			continue
		}
		goos := parts[0]
		goarch := parts[1]

		cmd := exec.Command("go", "build", "-v", "-o", output)
		cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH="+goarch)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("Build failed for %s with error: %v\n", platform, err)
			return
		}
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
