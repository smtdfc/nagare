package plugin

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/utils"
)

type Controller struct {
	pluginService *Service
}

func (p *Controller) List(ctx fiber.Ctx) error {
	data, err := p.pluginService.ListPlugins(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (p *Controller) InstallLocal(ctx fiber.Ctx) error {
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

func (p *Controller) Uninstall(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.UninstallPluginRequest](ctx)
	if err != nil {
		return err
	}

	err = p.pluginService.UninstallPlugin(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, struct{}{}, 200)
}

// @Injectable
func NewPluginController(pluginService *Service) *Controller {
	return &Controller{
		pluginService: pluginService,
	}
}
