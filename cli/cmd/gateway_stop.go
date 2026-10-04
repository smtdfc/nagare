package cmd

import (
	"fmt"

	"github.com/smtdfc/nagare/cli/helpers"
	"github.com/spf13/cobra"
)

var gatewayStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Nagare gateway service",
	Long:  "Request a graceful shutdown of the Nagare gateway using the PID recorded by gateway start.",
	Run: func(cmd *cobra.Command, args []string) {
		stopped, err := helpers.TryStopGateway()
		if err != nil {
			fmt.Printf("Failed to stop gateway: %v\n", err)
			return
		}

		if stopped {
			fmt.Println("Nagare gateway shutdown requested.")
			return
		}

		fmt.Println("Nagare gateway is not running.")
	},
}

func init() {
	gatewayCmd.AddCommand(gatewayStopCmd)
}
