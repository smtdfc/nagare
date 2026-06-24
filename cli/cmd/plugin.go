package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/smtdfc/nagare/core"
	"github.com/spf13/cobra"
)

// pluginCmd represents the plugin command
var installPluginCmd = &cobra.Command{
	Use:   "install [path_to_plugin]",
	Short: "Install a new plugin from a file",
	Long: `Install a Nagare plugin from a local .nagare_plugin archive file.

Example:
  nagare plugin install ./my-plugin.nagare_plugin
  nagare plugin install ~/downloads/data-connector.nagare_plugin`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pluginPath := args[0]
		pluginAbsPath, err := filepath.Abs(pluginPath)
		if err != nil {
			return fmt.Errorf("Error: %w", err)
		}

		err = core.PluginMgr.Install(pluginAbsPath)
		if err != nil {
			return err
		}

		fmt.Println("Plugin installed successfully")
		return nil
	},
}

// pluginCmd represents the plugin command
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage Nagare plugins",
	Long: `Manage, install, list, and remove plugins for the Nagare platform.
Plugins allow you to extend the functionality of the core system.
Use the subcommands below to perform specific actions on your plugins.`,
}

func init() {
	pluginCmd.AddCommand(installPluginCmd)
	rootCmd.AddCommand(pluginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pluginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pluginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
