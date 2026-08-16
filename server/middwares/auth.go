package middwares

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/custom_errors"
	"github.com/smtdfc/nagare/shared/helpers"
)

func AuthMiddleware(pubKey string) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		token := c.Get("X-Nagare-Secure")
		if token == "" {
			return custom_errors.NewServiceError("failed to authenticate", 403)
		}

		payload, err := helpers.VerifyToken(pubKey, token)
		if err != nil {
			fmt.Println("err", err)
			return custom_errors.NewServiceError("unauthorized", 403)
		}

		c.Locals("auth", payload)
		return c.Next()
	}
}
