package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/services"
	"github.com/smtdfc/nagare/gateway/utils"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

type ChatController struct {
	chatService *services.ChatService
}

func (c *ChatController) SendMessage(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.ChatSendMessageRequest](ctx)
	if err != nil {
		return err
	}

	err = c.chatService.SendMessage(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, 0, 200)
}

func (c *ChatController) CreateSession(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.ChatCreateSessionRequest](ctx)
	if err != nil {
		return err
	}

	data, err := c.chatService.CreateSession(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *ChatController) GetListSession(ctx fiber.Ctx) error {
	data, err := c.chatService.GetListSession(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *ChatController) GetHistory(ctx fiber.Ctx) error {
	sessionID := ctx.Query("session")

	data, err := c.chatService.GetHistory(ctx, sessionID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

// @Injectable
func NewChatController(chatService *services.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}
