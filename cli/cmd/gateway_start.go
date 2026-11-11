package cmd

import (
	"fmt"

	"github.com/smtdfc/nagare/cli/helpers"
	"github.com/spf13/cobra"
)

var debugMode bool
var gatewayBackground bool

var gatewayStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Nagare gateway service",
	Long:  "Start the Nagare gateway service. Use --debug to enable verbose logs and diagnostic information, and --background to detach it from this terminal.",

	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")

		if host != "" {
			fmt.Printf("Using host: %s\n", host)
		}

		err := helpers.TryStartGateway(debugMode, host, gatewayBackground)
		if err != nil {
			fmt.Printf("Failed to start gateway: %v\n", err)
		}
	},
}

func init() {
	gatewayCmd.AddCommand(gatewayStartCmd)

	gatewayStartCmd.Flags().BoolVarP(
		&debugMode,
		"debug",
		"d",
		false,
		"Enable verbose debug logging",
	)

	gatewayStartCmd.Flags().BoolVarP(
		&gatewayBackground,
		"background",
		"b",
		false,
		"Run the gateway in the background (detached from this terminal)",
	)

	gatewayStartCmd.Flags().String(
		"host",
		"",
		"Gateway host",
	)
}
