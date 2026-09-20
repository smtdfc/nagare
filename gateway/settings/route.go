package settings

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/guards"
	"github.com/smtdfc/nagare/gateway/common/middlewares"
	"github.com/smtdfc/nagare/gateway/common/websocket"
)

type SettingsRouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	settingsController *SettingsController,
	appConfig *config.Config,
	authGuard *guards.AuthGuard,
) SettingsRouteInitializer {
	authMiddleware := middlewares.AuthMiddlewareProvider(appConfig, authGuard)

	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Get(rest.GetGeneralSettings, authMiddleware, settingsController.GetGeneralConfig)
		app.Post(rest.SetGeneralSettings, authMiddleware, settingsController.SetGeneralConfig)
	}
}
