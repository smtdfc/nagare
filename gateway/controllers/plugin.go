package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/services"
	"github.com/smtdfc/nagare/gateway/utils"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

type PluginController struct {
	pluginService *services.PluginService
}

func (p *PluginController) GetListPlugin(ctx fiber.Ctx) error {
	data, err := p.pluginService.GetListPlugin(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (p *PluginController) InstallLocalPlugin(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.InstallLocalPluginRequest](ctx)
	if err != nil {
		return err
	}

	err = p.pluginService.InstallLocalPlugin(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, struct{}{}, 200)
}

// @Injectable
func NewPluginController(pluginService *services.PluginService) *PluginController {
	return &PluginController{
		pluginService: pluginService,
	}
}
