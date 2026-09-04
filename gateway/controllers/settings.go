package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/services"
	"github.com/smtdfc/nagare/gateway/utils"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

type SettingsController struct {
	settingsService *services.SettingsService
}

func (s *SettingsController) GetGeneralConfig(ctx fiber.Ctx) error {
	data, err := s.settingsService.GetGeneralSettings(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (s *SettingsController) SetGeneralConfig(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.SetGeneralSettingsRequest](ctx)
	if err != nil {
		return err
	}

	err = s.settingsService.SetGeneralSettings(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, 0, 200)
}

// @Injectable
func NewSettingsController(settingsService *services.SettingsService) *SettingsController {
	return &SettingsController{
		settingsService: settingsService,
	}
}
