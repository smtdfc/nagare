package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/desktop/security"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var isServerRunning = false

type AppService struct {
	id      string
	ctx     context.Context
	keypair *helpers.RSAKeys
}

func NewApp() *AppService {
	return &AppService{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *AppService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// app := application.Get()

	a.id = uuid.New().String()
	a.ctx = ctx
	// keypair, err := security.GetKeyPair()
	// if err != nil {
	// 	return nil
	// }
	// a.keypair = keypair

	isStart, err := helpers.CheckServerRun()
	if err != nil {
		fmt.Println(err)
	}

	// if !isStart {
	// 	mode := os.Getenv("NAGARE_MODE")
	// 	debugMode := false
	// 	if mode == "debug" {
	// 		debugMode = true
	// 	}

	// 	fmt.Printf("Server mode: %s\n", mode)
	// 	fmt.Println("Trying start server")

	// 	go func() {
	// 		err := helpers.TryStartServer(a.keypair.PublicKeyPEM, debugMode)
	// 		if err != nil {
	// 			fmt.Printf("Server runtime error: %v\n", err)

	// 		}
	// 	}()
	// }

	var serverReady bool
	for i := 1; i <= 5; i++ {
		fmt.Printf("Checking server health (attempt %d/5)...\n", i)

		if i > 1 || !isStart {
			time.Sleep(3 * time.Second)
		}

		running, checkErr := helpers.CheckServerRun()
		if checkErr == nil && running {
			serverReady = true
			fmt.Println("Server is up and running!")
			break
		}
	}

	if !serverReady {
		fmt.Println("Error: Server failed to start after multiple attempts.")
		// app.Dialog.Warning().
		// 	SetTitle("Error").
		// 	SetMessage(fmt.Sprintf("Server runtime error: %v\n", err)).
		// 	Show()
		// application.Get().Quit()
		return nil
	}
	isServerRunning = true
	return nil
}

func (a *AppService) ShowErrorDialog(title, message string) {
	application.Get().Dialog.Error().
		SetTitle(title).
		SetMessage(message).
		Show()
}

func (a *AppService) IsServerRunning() bool {
	return isServerRunning
}

func (a *AppService) GetWebsocketConnect() (string, error) {
	return helpers.GetWebsocketConnect()
}

func (a *AppService) GetRestApiConnect() (string, error) {
	return helpers.GetRestApiConnect()
}

func (a *AppService) GetToken() (string, error) {
	return security.GetToken()
}

func (a *AppService) SaveToken(token string) error {
	return security.SaveToken(token)
}

func (a *AppService) GetHost() (string, error) {
	return helpers.ResolveServerHost()
}

func (a *AppService) OpenPluginSelectDialog() string {
	path, err := application.Get().Dialog.OpenFile().
		SetTitle("Select PLugin file").
		AddFilter("Nagare Plugin", "*.nagare_plugin").
		AddFilter("All Files", "*.*").
		PromptForSingleSelection()

	if err != nil || path == "" {
		return ""
	}

	return path
}

func (a *AppService) ClearToken() error {
	return security.ClearToken()
}
