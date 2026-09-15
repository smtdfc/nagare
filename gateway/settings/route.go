package settings

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/middlewares"
	"github.com/smtdfc/nagare/gateway/common/websocket"
)

type RouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	settingsController *Controller,
	appConfig *config.Config,
) RouteInitializer {
	authMiddleware := middlewares.AuthMiddlewareProvider(appConfig)

	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Get(rest.GetGeneralSettings, authMiddleware, settingsController.GetGeneralConfig)
		app.Post(rest.SetGeneralSettings, authMiddleware, settingsController.SetGeneralConfig)
	}
}
