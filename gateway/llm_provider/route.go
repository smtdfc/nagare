package llm_provider

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/common/websocket"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

type RouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	llmProviderController *Controller,
) RouteInitializer {
	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Get(rest.ListLLMProvidersEndpoint, llmProviderController.List)
		app.Get(rest.GetLLMProviderDetailsEndpoint, llmProviderController.Details)
		app.Post(rest.AddLLMProviderEndpoint, llmProviderController.Add)
		app.Post(rest.DeleteLLMProviderEndpoint, llmProviderController.Delete)
		app.Post(rest.GetLLMProviderModelsEndpoint, llmProviderController.GetModels)
	}
}
