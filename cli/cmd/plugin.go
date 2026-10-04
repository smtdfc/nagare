package cmd

import (
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Build Nagare plugins",
	Long:  "Create, validate, and package plugins for the Nagare gateway.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
}
