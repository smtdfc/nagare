package credential

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
	credentialController *Controller,
	appConfig *config.Config,
	authGuard *guards.AuthGuard,
) RouteInitializer {
	authMiddleware := middlewares.AuthMiddlewareProvider(appConfig, authGuard)

	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Get(rest.ListCredentialsEndpoint, authMiddleware, credentialController.List)
		app.Get(rest.GetCredentialDetailsEndpoint, authMiddleware, credentialController.Details)
		app.Post(rest.AddCredentialEndpoint, authMiddleware, credentialController.Add)
		app.Post(rest.UpdateCredentialEndpoint, authMiddleware, credentialController.Update)
		app.Post(rest.DeleteCredentialEndpoint, authMiddleware, credentialController.Delete)
	}
}
