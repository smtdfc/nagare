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
	var fiberErr *fiber.Error
	var apiErr *rest.ApiError
	var nagareCoreErr *core_errors.NagareCoreError
	var gatewayErr *custom_errors.GatewayError
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
		errorResp = rest.NewApiError("ERROR", message, code)
	} else if errors.As(err, &apiErr) {
		errorResp = apiErr
	} else if errors.As(err, &nagareCoreErr) {
		errorResp = rest.NewApiError(
			nagareCoreErr.Code,
			nagareCoreErr.Details,
			400,
		)
	} else if errors.As(err, &gatewayErr) {
		errorResp = rest.NewApiError(
			gatewayErr.Code,
			gatewayErr.Details,
			gatewayErr.StatusCode,
		)
	}

	return utils.ResponseError(c, errorResp)
}
