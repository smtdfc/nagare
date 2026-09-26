package chat

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/common/custom_errors"
	"github.com/smtdfc/nagare/gateway/utils"
	"github.com/smtdfc/nagare/shared/security"
)

type ChatController struct {
	chatService *ChatService
	logger      *logger.BaseLogger
}

func authenticatedUserID(ctx fiber.Ctx) (string, error) {
	auth, ok := ctx.Locals("user").(*security.AuthPayload)
	if !ok || auth == nil || auth.ID == "" {
		return "", custom_errors.ErrUnauthorized
	}

	return auth.ID, nil
}

func (c *ChatController) SendMessage(ctx fiber.Ctx) error {
	ownerID, err := authenticatedUserID(ctx)
	if err != nil {
		return err
	}

	request, err := utils.ParseBody[*rest.SendChatMessageRequest](ctx)
	if err != nil {
		return err
	}

	err = c.chatService.SendMessage(ctx, ownerID, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, 0, 200)
}

func (c *ChatController) CreateSession(ctx fiber.Ctx) error {
	ownerID, err := authenticatedUserID(ctx)
	if err != nil {
		return err
	}

	request, err := utils.ParseBody[*rest.CreateChatSessionRequest](ctx)
	if err != nil {
		return err
	}

	data, err := c.chatService.CreateSession(ctx, ownerID, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *ChatController) ListSessions(ctx fiber.Ctx) error {
	ownerID, err := authenticatedUserID(ctx)
	if err != nil {
		return err
	}

	limit, _ := strconv.Atoi(ctx.Query("limit"))
	offset, _ := strconv.Atoi(ctx.Query("offset"))

	data, err := c.chatService.ListSessions(ctx, ownerID, offset, limit)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *ChatController) History(ctx fiber.Ctx) error {
	ownerID, err := authenticatedUserID(ctx)
	if err != nil {
		return err
	}

	sessionID := ctx.Params("id")
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	beforeID := ctx.Query("beforeID")

	data, err := c.chatService.GetHistoryPage(ctx, ownerID, sessionID, beforeID, limit)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *ChatController) GetSession(ctx fiber.Ctx) error {
	ownerID, err := authenticatedUserID(ctx)
	if err != nil {
		return err
	}

	sessionID := ctx.Params("id")
	data, err := c.chatService.GetSession(ctx, ownerID, sessionID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

// @Injectable
func NewController(chatService *ChatService, logger *logger.BaseLogger) *ChatController {
	return &ChatController{
		chatService: chatService,
		logger:      logger.With("module", "gateway:chat:controller"),
	}
}
