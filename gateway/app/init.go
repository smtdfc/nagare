package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
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
func StartApp(app *App, coreSetup *setup.CoreSetup, _ *Routes, chunkWorker *chat.Worker, _ *manager.PluginManager) *RunStats {
	err := coreSetup.Setup()
	if err != nil {
		return &RunStats{
			IsError: true,
			Error:   err,
		}
	}

	// The core chat worker handles SendEvent and runs the agent. This gateway
	// worker is a different subscriber: it forwards ChunkEvent payloads to the
	// session WebSocket rooms so clients receive the response as it streams.
	chunkWorker.Do()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		addr := fmt.Sprintf("%s:%s", app.config.Host, app.config.Port)
		err := app.fiberApp.Listen(addr, fiber.ListenConfig{
			// Replace Fiber's ASCII banner and process stats with a single
			// structured log line from the gateway logger.
			DisableStartupMessage: true,
			BeforeServeFunc: func(*fiber.App) error {
				app.logger.Info("Gateway started", "addr", addr)
				return nil
			},
		})
		if err != nil {
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

		err := coreSetup.Teardown(ctx)
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
