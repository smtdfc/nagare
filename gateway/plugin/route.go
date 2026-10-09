package plugin

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/guards"
	"github.com/smtdfc/nagare/gateway/common/middlewares"
	"github.com/smtdfc/nagare/gateway/common/websocket"
)

type RouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	pluginController *Controller,
	appConfig *config.Config,
	authGuard *guards.AuthGuard,
) RouteInitializer {
	authMiddleware := middlewares.AuthMiddlewareProvider(appConfig, authGuard)

	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Get(rest.GetListPluginEndpoint, authMiddleware, pluginController.List)
		app.Post(rest.InstallLocalPluginEndpoint, authMiddleware, pluginController.InstallLocal)
		app.Post(rest.UninstallPluginEndpoint, authMiddleware, pluginController.Uninstall)
		app.Post(rest.ActivatePluginEndpoint, authMiddleware, pluginController.Activate)
		app.Post(rest.DeactivatePluginEndpoint, authMiddleware, pluginController.Deactivate)
		app.Post(rest.GetPluginStatusEndpoint, authMiddleware, pluginController.Status)
		app.Post(rest.UploadPluginEndpoint, authMiddleware, pluginController.Upload)
		app.Post(rest.InstallPluginFromAttachmentEndpoint, authMiddleware, pluginController.InstallFromAttachment)
	}
}
