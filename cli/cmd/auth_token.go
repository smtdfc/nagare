package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/cli/security"
	"github.com/spf13/cobra"
)

var authTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		id := uuid.New().String()
		token, err := security.GenerateAuthToken(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Your token:", token)
	},
}

func init() {
	authCmd.AddCommand(authTokenCmd)
}
