package chat

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/utils"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

type Controller struct {
	chatService *Service
}

func (c *Controller) SendMessage(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.SendChatMessageRequest](ctx)
	if err != nil {
		return err
	}

	err = c.chatService.SendMessage(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, 0, 200)
}

func (c *Controller) CreateSession(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.CreateChatSessionRequest](ctx)
	if err != nil {
		return err
	}

	data, err := c.chatService.CreateSession(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *Controller) ListSessions(ctx fiber.Ctx) error {
	data, err := c.chatService.ListSessions(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *Controller) History(ctx fiber.Ctx) error {
	sessionID := ctx.Query("session")

	data, err := c.chatService.GetHistory(ctx, sessionID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

// @Injectable
func NewController(chatService *Service) *Controller {
	return &Controller{
		chatService: chatService,
	}
}
