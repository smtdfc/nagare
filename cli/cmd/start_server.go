package cmd

import (
	"fmt"
	"os"

	"github.com/smtdfc/nagare/cli/security"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/spf13/cobra"
)

var startServerCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		keypair, err := security.GetKeyPair()
		if err != nil {
			fmt.Println(err)
			return
		}

		mode := os.Getenv("NAGARE_MODE")
		debugMode := false
		if mode == "debug" {
			debugMode = true
		}

		err = helpers.TryStartServer(keypair.PublicKeyPEM, debugMode)
		if err != nil {
			fmt.Printf("Server runtime error: %v\n", err)
		}
	},
}

func init() {
	serverCmd.AddCommand(startServerCmd)
}
