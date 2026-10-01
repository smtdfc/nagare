package plugin

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/utils"
	"github.com/smtdfc/nagare/pkgs/paths"
)

type PluginController struct {
	logger        *logger.BaseLogger
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

	data, err := p.pluginService.InstallLocalPlugin(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
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

func (p *PluginController) Upload(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		p.logger.Error("Failed to get file from form file", "err", err)
		return custom_errors.ErrPluginFileInvalid
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".nagare_plugin" {
		return custom_errors.ErrPluginFileInvalidFormat
	}

	dir := filepath.Join(paths.UploadDir, "gateway/plugins")
	if err := os.MkdirAll(dir, 0755); err != nil {
		p.logger.Error("Failed to get file from form file", "err", err)
		return custom_errors.ErrUploadFailed
	}

	savePath := filepath.Join(dir, file.Filename)
	if err := c.SaveFile(file, savePath); err != nil {
		p.logger.Error("Failed to get file from form file", "err", err)
		return custom_errors.ErrUploadFailed
	}

	data, err := p.pluginService.UploadPlugin(c.Context(), savePath)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(c, data, 200)
}

func (p *PluginController) InstallFromAttachment(c fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.InstallPluginFromAttachmentRequest](c)
	if err != nil {
		return err
	}

	data, err := p.pluginService.InstallPluginFromAttachment(c.Context(), request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(c, data, 200)
}

// @Injectable
func NewPluginController(pluginService *PluginService, logger *logger.BaseLogger) *PluginController {
	return &PluginController{
		pluginService: pluginService,
		logger:        logger.With("module", "plugin-controller"),
	}
}
