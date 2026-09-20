package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/olahol/melody"
	config2 "github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/websocket"
)

type App struct {
	fiberApp      *fiber.App
	melody        *melody.Melody
	config        *config2.Config
	wsCoordinator *websocket.Coordinator
}

// @Injectable
func NewApp(
	config *config2.Config,
	wsCoordinator *websocket.Coordinator,
) *App {
	m := melody.New()
	m.Config.MaxMessageSize = 10 * 1024 * 1024 // 10MB

	return &App{
		fiberApp: fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		}),
		melody:        m,
		config:        config,
		wsCoordinator: wsCoordinator,
	}
}
