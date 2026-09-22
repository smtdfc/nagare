package app

import (
	"net/http"

	"github.com/gofiber/contrib/v3/monitor"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/smtdfc/nagare/gateway/auth"
	"github.com/smtdfc/nagare/gateway/chat"
	"github.com/smtdfc/nagare/gateway/llm_provider"
	"github.com/smtdfc/nagare/gateway/plugin"
	"github.com/smtdfc/nagare/gateway/settings"
)

type Routes struct{}

func RegisterPprofRoutes(app *fiber.App) {
	pprofGroup := app.Group("/debug/pprof")
	pprofGroup.Get("/*", adaptor.HTTPHandler(http.DefaultServeMux))
}

// @Injectable
func SetupRoutes(
	app *App,
	chatRoutes chat.ChatRouteInitializer,
	authRoutes auth.AuthRouteInitializer,
	llmProviderRoutes llm_provider.LLMProviderRouteInitializer,
	pluginRoutes plugin.PluginRouteInitializer,
	settingsRoutes settings.SettingsRouteInitializer,
) *Routes {

	chatRoutes(app.fiberApp, app.wsCoordinator)
	llmProviderRoutes(app.fiberApp, app.wsCoordinator)
	pluginRoutes(app.fiberApp, app.wsCoordinator)
	settingsRoutes(app.fiberApp, app.wsCoordinator)
	authRoutes(app.fiberApp, app.wsCoordinator)

	if app.config.DebugMode {
		app.fiberApp.Get("/metrics", monitor.New(monitor.Config{Title: "Nagare Gateway Metrics Page"}))
		RegisterPprofRoutes(app.fiberApp)
	}

	app.fiberApp.Use("/ws", func(c fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.fiberApp.Get("/ws", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := app.melody.HandleRequest(w, r)
		if err != nil {
			return
		}
	}))

	app.melody.HandleMessage(app.wsCoordinator.HandleMessage)
	app.melody.HandleDisconnect(app.wsCoordinator.HandleDisconnect)
	return &Routes{}
}
