package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/services"
	"github.com/smtdfc/nagare/server/utils"
	"github.com/smtdfc/nagare/shared/dto"
)

type SettingsController struct {
	service *services.SettingsService
}

func (c *SettingsController) GetGeneralSettings(ctx fiber.Ctx) error {
	resp, err := c.service.GetGeneralSettings()
	if err != nil {
		return err
	}

	return utils.SuccessResponse(resp, 200, ctx)
}

func (c *SettingsController) SaveGeneralSettings(ctx fiber.Ctx) error {
	var request dto.SaveGeneralSettingsRequest
	if err := ctx.Bind().JSON(&request); err != nil {
		return err
	}

	resp, err := c.service.SaveGeneralSettings(&request)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(resp, 200, ctx)
}

// @Injectable
func NewSettingsController(service *services.SettingsService) *SettingsController {
	return &SettingsController{service: service}
}
