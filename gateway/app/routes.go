package app

import (
	"net/http"

	"github.com/gofiber/contrib/v3/monitor"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/smtdfc/nagare/gateway/controllers"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

type AppRoutes struct{}

func RegisterPprofRoutes(app *fiber.App) {
	pprofGroup := app.Group("/debug/pprof")
	pprofGroup.Get("/*", adaptor.HTTPHandler(http.DefaultServeMux))
}

// @Injectable
func SetupRoutes(
	app *App,
	chatController *controllers.ChatController,
	llmProviderController *controllers.LLMProviderController,
	settingsController *controllers.SettingsController,
	pluginController *controllers.PluginController,
) *AppRoutes {
	app.fiberApp.Post(rest.SendChatMessageEndpoint, chatController.SendMessage)
	app.fiberApp.Post(rest.CreateChatSessionEndpoint, chatController.CreateSession)
	app.fiberApp.Get(rest.ListChatSessionsEndpoint, chatController.ListSessions)
	app.fiberApp.Get(rest.GetChatHistoryEndpoint, chatController.History)

	app.fiberApp.Get(rest.ListLLMProvidersEndpoint, llmProviderController.List)
	app.fiberApp.Get(rest.GetLLMProviderDetailsEndpoint, llmProviderController.Details)
	app.fiberApp.Post(rest.AddLLMProviderEndpoint, llmProviderController.Add)
	app.fiberApp.Post(rest.DeleteLLMProviderEndpoint, llmProviderController.Delete)
	app.fiberApp.Post(rest.GetLLMProviderModelsEndpoint, llmProviderController.GetModels)

	app.fiberApp.Get(rest.GetGeneralSettings, settingsController.GetGeneralConfig)
	app.fiberApp.Post(rest.SetGeneralSettings, settingsController.SetGeneralConfig)

	app.fiberApp.Get(rest.GetListPluginEndpoint, pluginController.List)
	app.fiberApp.Post(rest.InstallLocalPluginEndpoint, pluginController.InstallLocal)

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
	return &AppRoutes{}
}
