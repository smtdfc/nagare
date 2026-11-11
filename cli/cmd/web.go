package cmd

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/smtdfc/nagare/cli/helpers"
	"github.com/smtdfc/nagare/pkgs/system"
	"github.com/spf13/cobra"
)

var (
	webAddr       string
	webGateway    string
	webToken      string
	webNoOpen     bool
	webBackground bool
)

func defaultWebAddr() string {
	port := strings.TrimSpace(os.Getenv("NAGARE_WEB_UI_PORT"))
	if port == "" {
		port = "3005"
	}

	return "127.0.0.1:" + port
}

func defaultGatewayURL() string {
	if gatewayURL := strings.TrimSpace(os.Getenv("NAGARE_GATEWAY_URL")); gatewayURL != "" {
		return gatewayURL
	}

	port := strings.TrimSpace(os.Getenv("NAGARE_GATEWAY_PORT"))
	if port == "" {
		port = "9832"
	}

	return "http://localhost:" + port
}

func browserAddr(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}

	switch host {
	case "", "0.0.0.0", "::":
		host = "127.0.0.1"
	}

	return net.JoinHostPort(host, port)
}

func bootstrapURL(addr, gateway, token string) string {
	u := url.URL{
		Scheme: "http",
		Host:   browserAddr(addr),
		Path:   "/",
	}

	query := url.Values{}
	query.Set("gateway", gateway)
	if token != "" {
		query.Set("token", token)
	}
	u.RawQuery = query.Encode()

	return u.String()
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the Nagare Web UI",
	Long:  "Start the embedded Nagare Web UI and serve it over HTTP. A token is generated automatically and handed to the UI so the dashboard opens already signed in.",
	Run: func(cmd *cobra.Command, args []string) {
		if webBackground && !helpers.IsWebDaemon() {
			if err := helpers.StartBackgroundWebUI(browserAddr(webAddr)); err != nil {
				fmt.Printf("Failed to start Web UI in the background: %v\n", err)
			}
			return
		}

		if helpers.IsWebDaemon() {
			// The daemon owns the PID file; remove it when the process is
			// asked to terminate so `nagare web stop` stays accurate.
			helpers.InstallWebDaemonSignalCleanup()
		}

		token := strings.TrimSpace(webToken)
		if token == "" {
			generated, err := helpers.GenerateUserToken(helpers.DefaultTokenTTL)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Warning: could not generate a sign-in token: %v\n", err)
			} else {
				token = generated
			}
		}

		baseURL := fmt.Sprintf("http://%s", browserAddr(webAddr))
		signInURL := bootstrapURL(webAddr, webGateway, token)

		fmt.Printf("Nagare Web UI: %s\n", baseURL)
		if token != "" {
			fmt.Printf("Auto sign-in enabled for gateway %s\n", webGateway)
		}

		err := helpers.StartServer(webAddr, func() {
			if webNoOpen {
				fmt.Printf("Open this URL to sign in: %s\n", signInURL)
				return
			}

			if openErr := system.OpenBrowser(signInURL); openErr != nil {
				fmt.Printf("Open this URL to sign in: %s\n", signInURL)
			}
		})
		if err != nil {
			fmt.Printf("Failed to start Web UI: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVar(&webAddr, "addr", defaultWebAddr(), "HTTP address for the Web UI")
	webCmd.Flags().StringVar(&webGateway, "gateway", defaultGatewayURL(), "Gateway URL used for automatic sign-in")
	webCmd.Flags().StringVar(&webToken, "token", "", "Token used for automatic sign-in (a fresh token is generated when empty)")
	webCmd.Flags().BoolVar(&webNoOpen, "no-open", false, "Do not open the Web UI in a browser")
	webCmd.Flags().BoolVarP(&webBackground, "background", "b", false, "Run the Web UI in the background (detached from this terminal)")
}
