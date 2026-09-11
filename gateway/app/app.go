package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/gateway/common/websocket"
)

type App struct {
	fiberApp      *fiber.App
	melody        *melody.Melody
	config        *Config
	wsCoordinator *websocket.Coordinator
}

// @Injectable
func NewApp(
	config *Config,
	wsCoordinator *websocket.Coordinator,
) *App {

	return &App{
		fiberApp: fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		}),
		melody:        melody.New(),
		config:        config,
		wsCoordinator: wsCoordinator,
	}
}
