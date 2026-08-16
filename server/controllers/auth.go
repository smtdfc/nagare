package controllers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/smtdfc/nagare/server/custom_errors"
	"github.com/smtdfc/nagare/server/utils"
	"github.com/smtdfc/nagare/shared/dto"
)

type AuthController struct {
}

func (c *AuthController) Me(ctx fiber.Ctx) error {
	authPayload := ctx.Locals("auth").(*dto.AuthPayload)
	if authPayload == nil {
		return utils.ErrorResponse(custom_errors.NewServiceError("Unauthorized", 401), ctx)
	}

	return utils.SuccessResponse(&dto.GetProfileResponse{Profile: &dto.Profile{
		ID:   authPayload.ID,
		Role: authPayload.Role,
	}}, 200, ctx)
}

// @Injectable
func NewAuthController() *AuthController {
	return &AuthController{}
}
