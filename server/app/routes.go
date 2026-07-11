package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/controllers"
)

type AppRoute struct{}

// @Injectable
func InitRoutes(app *fiber.App, heathController *controllers.HealthController) *AppRoute {
	app.Get("/api/v1/health/check", heathController.CheckHealth)

	return &AppRoute{}
}
