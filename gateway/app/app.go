package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/gateway/event"
	"github.com/smtdfc/nagare/gateway/websocket"
)

type App struct {
	fiberApp      *fiber.App
	melody        *melody.Melody
	config        *AppConfig
	busSystem     *event.AppEventBusSystem
	wsCoordinator *websocket.WebsocketCoordinator
}

// @Injectable
func NewApp(

	config *AppConfig,
	busSys *event.AppEventBusSystem,
	wsCoordinator *websocket.WebsocketCoordinator,
) *App {

	return &App{
		fiberApp: fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		}),
		melody:        melody.New(),
		config:        config,
		busSystem:     busSys,
		wsCoordinator: wsCoordinator,
	}
}
