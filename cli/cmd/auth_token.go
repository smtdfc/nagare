package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	cli_helpers "github.com/smtdfc/nagare/cli/helpers"
	"github.com/smtdfc/nagare/pkgs/security"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate an authentication token",
	Long:  "Generate a JWT signed with the RSA key pair stored in the system keyring. Create a new 4096-bit key pair when no keys exist.",
	Run: func(cmd *cobra.Command, args []string) {
		_, privateKey, err := cli_helpers.GetRSAKey()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error generating RSA token: %v\n", err)
			return
		}

		payload := security.AuthPayload{
			ID:         uuid.Nil.String(),
			TargetType: "user",
			Name:       "User",
			Scopes:     []string{},
		}

		token, err := security.GenerateRSAToken(
			payload,
			[]byte(privateKey), 60*time.Hour,
		)
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
