package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/smtdfc/nagare/server/config"
	"github.com/smtdfc/nagare/server/utils"
)

// @Injectable
func NewFiberApp(config *config.ServerConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return utils.ErrorResponse(err, c)
		},
	})

	app.Use(func(c fiber.Ctx) error {
		if c.Method() == fiber.MethodOptions {
			c.Set("Access-Control-Allow-Origin", c.Get("Origin"))
			c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Nagare-Secure")
			c.Set("Access-Control-Allow-Credentials", "true")
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Nagare-Secure"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	return app
}

type App struct {
	fiberApp *fiber.App
	config   *config.ServerConfig
}

func (a *App) StartServer() *AppError {
	go func() {
		if err := a.fiberApp.Listen(fmt.Sprintf("127.0.0.1:%s", a.config.Port)); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := a.fiberApp.Shutdown(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped gracefully")
	return nil
}

// @Injectable
func NewApp(fiberApp *fiber.App, config *config.ServerConfig, _ *AppRoute) *App {
	return &App{
		fiberApp: fiberApp,
		config:   config,
	}
}
