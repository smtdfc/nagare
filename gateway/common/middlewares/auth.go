package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/custom_errors"
	"github.com/smtdfc/nagare/gateway/common/guards"
)

type AuthMiddleware fiber.Handler

func AuthMiddlewareProvider(
	appConfig *config.Config,
	authGuard *guards.AuthGuard,
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authorizationHeader := ctx.Get("Authorization")
		if authorizationHeader == "" {
			return custom_errors.ErrUnauthorized
		}

		parts := strings.Split(authorizationHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return custom_errors.ErrUnauthorized
		}

		tokenString := parts[1]
		if tokenString == "" {
			return custom_errors.ErrUnauthorized
		}

		auth, err := authGuard.VerifyUserFromToken(tokenString)
		if err != nil {
			return err
		}

		ctx.Locals("user", auth)
		return ctx.Next()
	}
}
