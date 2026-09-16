package llm_provider

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
	llmProviderController *Controller,
	appConfig *config.Config,
	authGuard *guards.AuthGuard,
) RouteInitializer {
	authMiddleware := middlewares.AuthMiddlewareProvider(appConfig, authGuard)

	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Get(rest.ListLLMProvidersEndpoint, authMiddleware, llmProviderController.List)
		app.Get(rest.GetLLMProviderDetailsEndpoint, authMiddleware, llmProviderController.Details)
		app.Post(rest.AddLLMProviderEndpoint, authMiddleware, llmProviderController.Add)
		app.Post(rest.DeleteLLMProviderEndpoint, authMiddleware, llmProviderController.Delete)
		app.Post(rest.GetLLMProviderModelsEndpoint, authMiddleware, llmProviderController.GetModels)
	}
}
