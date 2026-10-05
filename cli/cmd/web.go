package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/smtdfc/nagare/cli/helpers"
	"github.com/spf13/cobra"
)

var webAddr string

func defaultWebAddr() string {
	port := strings.TrimSpace(os.Getenv("NAGARE_WEB_UI_PORT"))
	if port == "" {
		port = "3005"
	}

	return "127.0.0.1:" + port
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the Nagare Web UI",
	Long:  "Start the embedded Nagare Web UI and serve it over HTTP.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := helpers.StartServer(webAddr); err != nil {
			fmt.Printf("Failed to start Web UI: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVar(&webAddr, "addr", defaultWebAddr(), "HTTP address for the Web UI")
}
