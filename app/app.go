package main

import (
	"context"
	"fmt"

	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	isStart, err := helpers.CheckServerRun("3000")
	if err != nil {
		fmt.Println(err)
	}

	if !isStart {
		err = helpers.TryStartServer("3000")
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
	return helpers.GetWebsocketConnect("3000")
}
