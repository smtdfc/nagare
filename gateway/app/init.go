package app

import (
	"fmt"

	"github.com/smtdfc/nagare/core/setup"
	"github.com/smtdfc/nagare/gateway/chat"
)

type RunStats struct {
	IsError bool
	Error   error
}

// @Injectable
// @Root
func StartApp(app *App, coreSetup *setup.CoreSetup, _ *Routes, chatWorker *chat.Worker) *RunStats {
	err := coreSetup.Setup()
	if err != nil {
		return &RunStats{
			IsError: true,
			Error:   err,
		}
	}
	chatWorker.Start()
	err = app.fiberApp.Listen(fmt.Sprintf("localhost:%s", app.config.Port))
	if err != nil {
		return &RunStats{
			IsError: true,
			Error:   err,
		}
	}

	return &RunStats{
		IsError: false,
	}
}
