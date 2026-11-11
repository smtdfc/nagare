package cmd

import (
	"fmt"

	"github.com/smtdfc/nagare/cli/helpers"
	"github.com/spf13/cobra"
)

var webStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the background Nagare Web UI daemon",
	Long:  "Stop the background Nagare Web UI daemon started by 'nagare web -b'. Reports when no daemon is running.",

	Run: func(cmd *cobra.Command, args []string) {
		stopped, err := helpers.TryStopWebUI()
		switch {
		case err != nil:
			fmt.Printf("Failed to stop Web UI: %v\n", err)
		case stopped:
			fmt.Println("Nagare Web UI stopped")
		default:
			fmt.Println("No background Nagare Web UI is running")
		}
	},
}

func init() {
	webCmd.AddCommand(webStopCmd)
}
