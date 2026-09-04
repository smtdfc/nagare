package utils

import (
	"bytes"

	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/custom_errors"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

func ResponseSuccess[T any](ctx fiber.Ctx, data T, code int) error {
	return ctx.Status(code).JSON(&rest.ApiResponse[T]{
		IsSuccess: true,
		Data:      data,
	})
}

func ResponseError(ctx fiber.Ctx, err *rest.ApiError) error {
	return ctx.Status(err.StatusCode).JSON(&rest.ApiResponse[any]{
		IsSuccess: false,
		Error:     err,
	})
}

func ParseBody[T any](ctx fiber.Ctx) (T, error) {
	var data T

	body := bytes.TrimSpace(ctx.Body())
	if len(body) == 0 {
		return data, custom_errors.ErrMissingRequestBody
	}

	if err := ctx.Bind().Body(&data); err != nil {
		return data, err
	}

	return data, nil
}
