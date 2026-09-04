package cmd

import (
	"fmt"

	"github.com/smtdfc/nagare/cli/helpers"
	"github.com/spf13/cobra"
)

var debugMode bool

var gatewayStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Nagare gateway service",
	Long: `Start the Nagare gateway service. 
Use the --debug flag to run the gateway in debug mode, providing more verbose logs and diagnostic information.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := helpers.TryStartGateway(debugMode)
		if err != nil {
			fmt.Printf("Failed to start gateway: %v\n", err)
		}
	},
}

func init() {
	gatewayCmd.AddCommand(gatewayStartCmd)
	gatewayStartCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug mode for verbose logging")
}
