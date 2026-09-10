package app

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	core_errors "github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/gateway/custom_errors"
	"github.com/smtdfc/nagare/gateway/utils"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	errorResp := rest.InternalErr
	if fiberErr, ok := errors.AsType[*fiber.Error](err); ok {
		code = fiberErr.Code
		message = fiberErr.Message
		errorResp = rest.NewApiError("ERROR", message, code)
	} else if apiErr, ok := errors.AsType[*rest.ApiError](err); ok {
		errorResp = apiErr
	} else if nagareCoreErr, ok := errors.AsType[*core_errors.NagareCoreError](err); ok {
		errorResp = rest.NewApiError(
			nagareCoreErr.Code,
			nagareCoreErr.Details,
			400,
		)
	} else if gatewayErr, ok := errors.AsType[*custom_errors.GatewayError](err); ok {
		errorResp = rest.NewApiError(
			gatewayErr.Code,
			gatewayErr.Details,
			gatewayErr.StatusCode,
		)
	}

	return utils.ResponseError(c, errorResp)
}
