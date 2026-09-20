package chat

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	websocket_dtos "github.com/smtdfc/nagare/dtos/websocket"
	"github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/guards"
	"github.com/smtdfc/nagare/gateway/common/middlewares"
	"github.com/smtdfc/nagare/gateway/common/websocket"
)

type ChatRouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	chatController *ChatController,
	websocketHandler *ChatWebsocketHandler,
	appConfig *config.Config,
	authGuard *guards.AuthGuard,
) ChatRouteInitializer {
	authMiddleware := middlewares.AuthMiddlewareProvider(appConfig, authGuard)

	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Post(rest.SendChatMessageEndpoint, authMiddleware, chatController.SendMessage)
		app.Post(rest.CreateChatSessionEndpoint, authMiddleware, chatController.CreateSession)
		app.Get(rest.ListChatSessionsEndpoint, authMiddleware, chatController.ListSessions)
		app.Get(rest.GetChatHistoryEndpoint, authMiddleware, chatController.History)

		ws.On(websocket_dtos.RegisterChatListenerEvent, websocketHandler.OnListenMessage)
		ws.On(websocket_dtos.AuthEvent, websocketHandler.OnAuth)
	}
}
