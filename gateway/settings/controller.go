package settings

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/utils"
)

type Controller struct {
	settingsService *SettingsService
}

func (s *Controller) GetGeneralConfig(ctx fiber.Ctx) error {
	data, err := s.settingsService.GetGeneralSettings(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (s *Controller) SetGeneralConfig(ctx fiber.Ctx) error {
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
func NewSettingsController(settingsService *SettingsService) *Controller {
	return &Controller{
		settingsService: settingsService,
	}
}
