package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags "-X github.com/smtdfc/nagare/cli/cmd.Version=x.y.z".
// Falls back to "dev" when building without ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "nagare",
	Short: "Command-line tools for managing Nagare",
	Long:  "Nagare CLI provides commands for authentication, gateway management, and plugin development.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if v, _ := cmd.Flags().GetBool("version"); v {
			fmt.Printf("nagare version %s\n", Version)
			return nil
		}
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Accept the documented -sign form as an alias for --sign.
	rootCmd.SetArgs(normalizeLegacyFlags(os.Args[1:]))

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func normalizeLegacyFlags(args []string) []string {
	normalized := make([]string, len(args))
	for i, arg := range args {
		switch {
		case arg == "-sign":
			normalized[i] = "--sign"
		case strings.HasPrefix(arg, "-sign="):
			normalized[i] = "--" + strings.TrimPrefix(arg, "-")
		default:
			normalized[i] = arg
		}
	}
	return normalized
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.Flags().BoolP("version", "v", false, "Print version information and exit")
}
