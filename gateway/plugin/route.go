package plugin

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/common/websocket"
	plugin_dtos "github.com/smtdfc/nagare/shared/dtos/plugin"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

type RouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	pluginController *Controller,
	websocketHandler *WebsocketHandler,
) RouteInitializer {
	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Get(rest.GetListPluginEndpoint, pluginController.List)
		app.Post(rest.InstallLocalPluginEndpoint, pluginController.InstallLocal)

		ws.On(plugin_dtos.HandshakeEvent, websocketHandler.OnHandshakeEvent)
	}
}
