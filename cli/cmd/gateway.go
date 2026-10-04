package cmd

import (
	"github.com/spf13/cobra"
)

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Manage the Nagare gateway",
	Long:  "Start and manage the Nagare gateway service.",
}

func init() {
	rootCmd.AddCommand(gatewayCmd)
}
