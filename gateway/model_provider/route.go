package model_provider

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/guards"
	"github.com/smtdfc/nagare/gateway/common/middlewares"
	"github.com/smtdfc/nagare/gateway/common/websocket"
)

type ModelProviderRouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	modelProviderController *ModelProviderController,
	appConfig *config.Config,
	authGuard *guards.AuthGuard,
) ModelProviderRouteInitializer {
	authMiddleware := middlewares.AuthMiddlewareProvider(appConfig, authGuard)

	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Get(rest.ListModelProvidersEndpoint, authMiddleware, modelProviderController.List)
		app.Get(rest.GetModelProviderDetailsEndpoint, authMiddleware, modelProviderController.Details)
		app.Post(rest.AddModelProviderEndpoint, authMiddleware, modelProviderController.Add)
		app.Post(rest.DeleteModelProviderEndpoint, authMiddleware, modelProviderController.Delete)
		app.Post(rest.GetModelProviderModelsEndpoint, authMiddleware, modelProviderController.GetModels)
	}
}
