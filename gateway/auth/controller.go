package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/common/custom_errors"
	"github.com/smtdfc/nagare/gateway/utils"
)

type AuthController struct {
	logger *logger.BaseLogger
}

func (a *AuthController) Check(ctx fiber.Ctx) error {
	user := ctx.Locals("user")
	if user == nil {
		return custom_errors.ErrUnauthorized
	}

	return utils.ResponseSuccess(ctx, &rest.CheckAuthStateResponse{IsAuth: true}, 200)
}

// @Injectable
func NewAuthController(logger *logger.BaseLogger) *AuthController {
	return &AuthController{
		logger: logger.With("module", "auth-controller"),
	}
}
