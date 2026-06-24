package cmd

import (
	"os"

	"github.com/smtdfc/nagare/cli/tui"
	"github.com/smtdfc/nagare/core"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nagare",
	Short: "Nagare AI Agent",
	Long:  ``,

	// Disable default completion if not needed for cleaner UI
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := core.PreStart()
		if err != nil {
			return err
		}

		tui.NewRootTUI(core.Config)
		core.Shutdown()
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
