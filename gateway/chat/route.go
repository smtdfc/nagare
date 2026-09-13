package chat

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/middlewares"
	"github.com/smtdfc/nagare/gateway/common/websocket"
	"github.com/smtdfc/nagare/shared/dtos/rest"
	websocket_dtos "github.com/smtdfc/nagare/shared/dtos/websocket"
)

type RouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	chatController *Controller,
	websocketHandler *WebsocketHandler,
	appConfig *config.Config,
) RouteInitializer {
	authMiddleware := middlewares.AuthMiddlewareProvider(appConfig)

	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Post(rest.SendChatMessageEndpoint, authMiddleware, chatController.SendMessage)
		app.Post(rest.CreateChatSessionEndpoint, authMiddleware, chatController.CreateSession)
		app.Get(rest.ListChatSessionsEndpoint, authMiddleware, chatController.ListSessions)
		app.Get(rest.GetChatHistoryEndpoint, authMiddleware, chatController.History)

		ws.On(websocket_dtos.RegisterChatListenerEvent, websocketHandler.OnListenMessage)
	}
}
