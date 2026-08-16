package app

import (
	"log"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/config"
	"github.com/smtdfc/nagare/server/controllers"
	"github.com/smtdfc/nagare/server/middwares"
	"github.com/smtdfc/nagare/server/ws"
	"github.com/smtdfc/nagare/server/ws/handlers"
	"github.com/smtdfc/nagare/shared/dto"
)

type AppRoute struct{}

// @Injectable
func InitRoutes(
	app *fiber.App,
	heathController *controllers.HealthController,
	authController *controllers.AuthController,
	providerController *controllers.ProviderController,
	settingsController *controllers.SettingsController,
	pluginController *controllers.PluginController,
	chatHandler *handlers.ChatHandler,
	authHanlder *handlers.AuthHandler,
	config *config.ServerConfig,
) *AppRoute {
	app.Get("/api/v1/health/check", heathController.CheckHealth)
	app.Get("/api/v1/auth/me", middwares.AuthMiddleware(config.PublicKey), authController.Me)
	app.Get("/api/v1/provider/list", middwares.AuthMiddleware(config.PublicKey), providerController.GetListProvider)
	app.Get("/api/v1/provider/:id/details", middwares.AuthMiddleware(config.PublicKey), providerController.GetProviderDetails)
	app.Post("/api/v1/provider/create", middwares.AuthMiddleware(config.PublicKey), providerController.CreateProvider)
	app.Post("/api/v1/provider/update", middwares.AuthMiddleware(config.PublicKey), providerController.UpdateProvider)
	app.Post("/api/v1/provider/delete", middwares.AuthMiddleware(config.PublicKey), providerController.DeleteProvider)
	app.Post("/api/v1/provider/fetch-model", middwares.AuthMiddleware(config.PublicKey), providerController.FetchModel)
	app.Get("/api/v1/settings/general", settingsController.GetGeneralSettings)
	app.Post("/api/v1/settings/general/save", settingsController.SaveGeneralSettings)
	app.Get("/api/v1/plugin/list", pluginController.GetAll)
	app.Post("/api/v1/plugin/install-local", pluginController.InstallLocalPlugin)
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
			case dto.WS_INVOKE_AGENT:
				chatHandler.InvokeAgent(instance, wsMsg)
			case dto.WS_AUTH_REQUEST:
				authHanlder.Auth(instance, wsMsg)
			}
		}
	}))
	return &AppRoute{}
}
