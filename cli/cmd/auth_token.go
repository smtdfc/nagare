package cmd

import (
	"fmt"
	"os"

	"github.com/smtdfc/nagare/cli/helpers"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate an authentication token",
	Long:  "Generate a JWT signed with the RSA key pair stored in the system keyring. Create a new 4096-bit key pair when no keys exist.",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := helpers.GenerateUserToken(helpers.DefaultTokenTTL)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error generating RSA token: %v\n", err)
			return
		}

		fmt.Printf("Your token is: %s \n", token)
	},
}

func init() {
	authCmd.AddCommand(tokenCmd)

}
