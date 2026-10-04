package cmd

import (
	"github.com/smtdfc/nagare/plugin/builder"
	"github.com/spf13/cobra"
)

var pluginBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and package a Nagare plugin",
	Long:  "Validate plugin metadata, compile and sign the plugin binary, and package it as a .nagare_plugin archive.",
	Run: func(cmd *cobra.Command, args []string) {
		builder.Build()
	},
}

func init() {
	pluginCmd.AddCommand(pluginBuildCmd)
}
