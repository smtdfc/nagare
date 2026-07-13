package utils

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/custom_errors"
	"github.com/smtdfc/nagare/shared/dto"
)

func SuccessResponse[T any](data *T, status int, ctx fiber.Ctx) error {
	return ctx.Status(status).JSON(dto.ApiResponse[*T]{
		Status: "success",
		Data:   data,
	})
}

func ErrorResponse(err error, ctx fiber.Ctx) error {
	status := 500
	errResp := dto.ApiError{
		Name:    "",
		Message: "",
		Details: map[string]string{},
	}

	if errors.As(err, &custom_errors.ServiceError{}) {
		errResp.Name = "ServiceError"
		errResp.Message = err.Error()
		status = err.(custom_errors.ServiceError).Status
	} else {
		errResp.Name = "InternalServiceError"
		errResp.Message = "InternalServiceError"
	}

	return ctx.Status(status).JSON(dto.ApiResponse[any]{
		Status: "error",
		Error:  &errResp,
	})
}
