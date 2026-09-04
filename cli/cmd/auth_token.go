package cmd

import (
	"fmt"
	"os"
	"time"

	cli_helpers "github.com/smtdfc/nagare/cli/helpers"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/security"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate a token",
	Long: `Generate a secure JWT token signed with an RSA private key. 
This command automatically checks the system keyring for existing RSA keys. 
If no keys are found, it generates a new 4096-bit RSA key pair, stores them securely, 
and uses them to issue a new authentication token with a predefined user payload.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, privateKey, err := cli_helpers.GetRSAKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating RSA token: %v\n", err)
			return
		}

		payload := security.UserAuthPayload{
			ID:     helpers.GenerateUUID(),
			Name:   "User",
			Scopes: []string{},
		}

		token, err := security.GenerateRSAToken(
			payload,
			[]byte(privateKey), 60*time.Hour,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating RSA token: %v\n", err)
			return
		}

		fmt.Printf("Your token is: %s \n", token)
	},
}

func init() {
	authCmd.AddCommand(tokenCmd)

}
