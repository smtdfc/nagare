package plugin

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/utils"
)

type PluginController struct {
	pluginService *PluginService
}

func (p *PluginController) List(ctx fiber.Ctx) error {
	data, err := p.pluginService.ListPlugins(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (p *PluginController) InstallLocal(ctx fiber.Ctx) error {
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

func (p *PluginController) Uninstall(ctx fiber.Ctx) error {
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

func (p *PluginController) Activate(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.ActivatePluginRequest](ctx)
	if err != nil {
		return err
	}

	err = p.pluginService.ActivatePlugin(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, struct{}{}, 200)
}

func (p *PluginController) Deactivate(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.DeactivatePluginRequest](ctx)
	if err != nil {
		return err
	}

	err = p.pluginService.DeactivatePlugin(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, struct{}{}, 200)
}

func (p *PluginController) Status(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.GetPluginStatusRequest](ctx)
	if err != nil {
		return err
	}

	data, err := p.pluginService.GetPluginStatus(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

// @Injectable
func NewPluginController(pluginService *PluginService) *PluginController {
	return &PluginController{
		pluginService: pluginService,
	}
}
