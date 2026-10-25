package cmd

import (
	"github.com/smtdfc/nagare/plugin/builder"
	"github.com/spf13/cobra"
)

var pluginSignKeyPath string

var pluginBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and package a Nagare plugin",
	Long:  "Validate plugin metadata, compile the plugin binary, package it as a .nagare_plugin archive, and optionally sign the archive with an Ed25519 private key.",
	Run: func(cmd *cobra.Command, args []string) {
		builder.Build(pluginSignKeyPath)
	},
}

func init() {
	pluginBuildCmd.Flags().StringVarP(&pluginSignKeyPath, "sign", "s", "", "Path to an Ed25519 private key used to sign the packaged archive")
	pluginCmd.AddCommand(pluginBuildCmd)
}
