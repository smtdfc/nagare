package app

import (
	"net/http"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/smtdfc/nagare/server/config"
	"github.com/smtdfc/nagare/server/controllers"
	"github.com/smtdfc/nagare/server/middwares"
	"github.com/smtdfc/nagare/server/ws/handlers"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/ws"
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
	pluginHandler *handlers.PluginHandler,
	config *config.ServerConfig,
) *AppRoute {
	authMidleware := middwares.AuthMiddleware(config.PublicKey)
	app.Get("/api/v1/health/check", heathController.CheckHealth)
	app.Get("/api/v1/auth/me", authMidleware, authController.Me)
	app.Get("/api/v1/provider/list", authMidleware, providerController.GetListProvider)
	app.Get("/api/v1/provider/:id/details", authMidleware, providerController.GetProviderDetails)
	app.Post("/api/v1/provider/create", authMidleware, providerController.CreateProvider)
	app.Post("/api/v1/provider/update", authMidleware, providerController.UpdateProvider)
	app.Post("/api/v1/provider/delete", authMidleware, providerController.DeleteProvider)
	app.Post("/api/v1/provider/fetch-model", authMidleware, providerController.FetchModel)
	app.Get("/api/v1/settings/general", authMidleware, settingsController.GetGeneralSettings)
	app.Post("/api/v1/settings/general/save", authMidleware, settingsController.SaveGeneralSettings)
	app.Get("/api/v1/plugin/list", authMidleware, pluginController.GetAll)
	app.Post("/api/v1/plugin/install-local", authMidleware, pluginController.InstallLocalPlugin)
	app.Post("/api/v1/plugin/remove", authMidleware, pluginController.RemovePlugin)

	app.Use("/ws", func(c fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/chat", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := ws.WsHandler(w, r, func(instance *ws.WsInstance, wsMsg *dto.WsMessage[any]) error {
			switch wsMsg.Event {
			case dto.WS_CREATE_SESSION:
				chatHandler.CreateSession(instance)
			case dto.WS_INVOKE_AGENT:
				chatHandler.InvokeAgent(instance, wsMsg)
			case dto.WS_AUTH_REQUEST:
				authHanlder.Auth(instance, wsMsg)
			}
			return nil
		})

		if err != nil {
			println("WebSocket error:", err.Error())
		}
	}))

	app.Get("/ws/plugin", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := ws.WsHandler(w, r, func(instance *ws.WsInstance, wsMsg *dto.WsMessage[any]) error {
			switch wsMsg.Event {
			case dto.WS_PLUGIN_REGISTER:
				pluginHandler.Register(instance, wsMsg)
			}
			return nil
		})

		if err != nil {
			println("WebSocket error:", err.Error())
		}
	}))

	return &AppRoute{}
}
