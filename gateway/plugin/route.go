package plugin

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/common/websocket"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

type RouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	pluginController *Controller,
) RouteInitializer {
	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Get(rest.GetListPluginEndpoint, pluginController.List)
		app.Post(rest.InstallLocalPluginEndpoint, pluginController.InstallLocal)
	}
}
