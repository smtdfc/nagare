package cmd

import (
	"os"

	"github.com/smtdfc/nagare/cli/tui"
	"github.com/smtdfc/nagare/core/config"
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
		conf, err := config.LoadConfig()
		if err != nil {
			conf = &config.Config{
				CurrentMode: config.PROVIDER_MODE,
				Providers: map[string]config.Provider{
					"OpenAI": {
						Name:       "OpenAI",
						BaseURL:    "https://api.openai.com/v1",
						Compatible: config.OPEN_AI_COMPATIBLE,
						Enabled:    true,
					},
				},
			}
		}

		tui.NewRootTUI(conf)
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nagare.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
