package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/utils"
)

type AppError struct {
	error
}

// @Injectable
func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return utils.ErrorResponse(err, c)
		},
	})

	return app
}

// @Injectable
// @Root
func StartApp(app *fiber.App, _ *AppRoute) *AppError {
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped gracefully")
	return nil
}
