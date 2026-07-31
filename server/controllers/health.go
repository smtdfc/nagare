package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/services"
	"github.com/smtdfc/nagare/server/utils"
)

type HealthController struct {
	healthService *services.HealthService
}

func (h *HealthController) CheckHealth(ctx fiber.Ctx) error {
	resp, err := h.healthService.CheckHealth()
	if err != nil {
		return err
	}

	return utils.SuccessResponse(resp, 200, ctx)
}

// @Injectable
func NewHealthController(healthService *services.HealthService) *HealthController {
	return &HealthController{
		healthService: healthService,
	}
}
