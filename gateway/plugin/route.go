package plugin

import (
	"github.com/gofiber/fiber/v3"
	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/middlewares"
	"github.com/smtdfc/nagare/gateway/common/websocket"
)

type RouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	pluginController *Controller,
	websocketHandler *WebsocketHandler,
	appConfig *config.Config,
) RouteInitializer {
	authMiddleware := middlewares.AuthMiddlewareProvider(appConfig)
	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Get(rest.GetListPluginEndpoint, authMiddleware, pluginController.List)
		app.Post(rest.InstallLocalPluginEndpoint, authMiddleware, pluginController.InstallLocal)
		app.Post(rest.UninstallPluginEndpoint, authMiddleware, pluginController.Uninstall)
		app.Post(rest.ActivatePluginEndpoint, authMiddleware, pluginController.Activate)
		app.Post(rest.DeactivatePluginEndpoint, authMiddleware, pluginController.Deactivate)
		app.Post(rest.GetPluginStatusEndpoint, authMiddleware, pluginController.Status)
		ws.On(plugin_dtos.HandshakeEvent, websocketHandler.OnHandshakeEvent)
		ws.On(plugin_dtos.PrepareChatSessionEvent, websocketHandler.OnPrepareChatSession)
		ws.On(plugin_dtos.SendChatMessageEvent, websocketHandler.OnSendChatMessage)
	}
}
