package app

import (
	"log"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/controllers"
	"github.com/smtdfc/nagare/server/ws"
	"github.com/smtdfc/nagare/server/ws/handlers"
	"github.com/smtdfc/nagare/shared/dto"
)

type AppRoute struct{}

// @Injectable
func InitRoutes(app *fiber.App, heathController *controllers.HealthController, chatHandler *handlers.ChatHandler) *AppRoute {
	app.Get("/api/v1/health/check", heathController.CheckHealth)
	app.Use("/ws", func(c fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/chat", websocket.New(func(c *websocket.Conn) {
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		instance := ws.NewWsInstance(c)
		for {
			wsMsg, err := ws.ReadMessage[any](instance)
			if err != nil {
				log.Println("read err:", err)
				break
			}

			switch wsMsg.Event {
			case dto.WS_CREATE_SESSION:
				chatHandler.CreateSession(instance)
			}
		}
	}))
	return &AppRoute{}
}
