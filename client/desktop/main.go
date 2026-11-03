package main

import (
	"embed"

	"log"
	"net/url"
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[string]("time")
}

// bootstrapURL forwards launcher-provided connection settings to the frontend.
// The window URL format matches `nagare web`, so the shared UI can store the
// token and gateway and open already signed in.
func bootstrapURL() string {
	query := url.Values{}
	if gateway := strings.TrimSpace(os.Getenv("NAGARE_GATEWAY_URL")); gateway != "" {
		query.Set("gateway", gateway)
	}
	if token := strings.TrimSpace(os.Getenv("NAGARE_AUTH_TOKEN")); token != "" {
		query.Set("token", token)
	}

	if len(query) == 0 {
		return "/"
	}

	return "/?" + query.Encode()
}

func main() {

	app := application.New(application.Options{
		Name:        "nagare-desktop",
		Description: "",
		Services:    []application.Service{},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Nagare Desktop",
		Width:  1200,
		Height: 700,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		URL: bootstrapURL(),
	})

	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
