package chat

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/common/websocket"
	"github.com/smtdfc/nagare/shared/dtos/rest"
	websocket_dtos "github.com/smtdfc/nagare/shared/dtos/websocket"
)

type RouteInitializer func(app *fiber.App, ws *websocket.Coordinator)

// @Injectable
func NewRouteInitializer(
	chatController *Controller,
	websocketHandler *WebsocketHandler,
) RouteInitializer {
	return func(app *fiber.App, ws *websocket.Coordinator) {
		app.Post(rest.SendChatMessageEndpoint, chatController.SendMessage)
		app.Post(rest.CreateChatSessionEndpoint, chatController.CreateSession)
		app.Get(rest.ListChatSessionsEndpoint, chatController.ListSessions)
		app.Get(rest.GetChatHistoryEndpoint, chatController.History)

		ws.On(websocket_dtos.ReceivedChatMessageEvent, websocketHandler.OnListenMessage)
	}
}
