package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/core/setup"
	"github.com/smtdfc/nagare/gateway/chat"
)

type RunStats struct {
	IsError bool
	Error   error
}

// @Injectable
// @Root
func StartApp(app *App, coreSetup *setup.CoreSetup, _ *Routes, chatWorker *chat.Worker, pluginMgr *manager.PluginManager) *RunStats {
	err := coreSetup.Setup(app.config.Port)
	if err != nil {
		return &RunStats{
			IsError: true,
			Error:   err,
		}
	}

	chatWorker.Start()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		addr := fmt.Sprintf("localhost:%s", app.config.Port)
		if err := app.fiberApp.Listen(addr); err != nil {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return &RunStats{
			IsError: true,
			Error:   err,
		}
	case sig := <-quit:
		log.Printf("Received stop signal (%v), initiating graceful shutdown...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := pluginMgr.StopAllPlugin(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		if err := app.fiberApp.ShutdownWithContext(ctx); err != nil {
			log.Printf("Error: %v", err)
		}
	}

	return &RunStats{
		IsError: false,
	}
}
