package main

import (
	"app/security"
	"context"
	"fmt"

	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	keypair *helpers.RSAKeys
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func ShowErrorDialog(ctx context.Context, message string) (string, error) {
	return runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   "Error",
		Message: message,
	})
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	keypair, err := security.GetKeyPair()
	if err != nil {
		ShowErrorDialog(ctx, "Keyring error!")
		return
	}
	a.keypair = keypair

	isStart, err := helpers.CheckServerRun("9832")
	if err != nil {
		fmt.Println(err)
	}

	if !isStart {
		fmt.Println("Trying start server")
		err = helpers.TryStartServer("9832", a.keypair.PublicKeyPEM)
	}

	if err != nil {
		fmt.Println(err)
		_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Error",
			Message: "Initialization error: Nagare server failed to start. Please check the configuration and logs.",
		})
		if err != nil {

		}

		runtime.Quit(a.ctx)
		return
	}
}

func (a *App) GetWebsocketConnect() (string, error) {
	return helpers.GetWebsocketConnect("9832")
}
