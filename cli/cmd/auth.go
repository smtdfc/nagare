package cmd

import (
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication credentials",
	Long:  "Create and manage the credentials used to authenticate Nagare services.",
}

func init() {
	rootCmd.AddCommand(authCmd)
}
