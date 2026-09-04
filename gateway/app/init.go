package app

import (
	"fmt"

	"github.com/smtdfc/nagare/core/setup"
	"github.com/smtdfc/nagare/gateway/workers"
)

type AppError struct {
	error
}

// @Injectable
// @Root
func StartApp(app *App, coreSetup *setup.CoreSetup, _ *AppRoutes, chatWorker *workers.ChatWorker) *AppError {
	coreSetup.Setup()
	chatWorker.Start()
	app.fiberApp.Listen(fmt.Sprintf("localhost:%s", app.config.Port))
	return nil
}
