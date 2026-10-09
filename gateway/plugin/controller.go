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

type Controller struct {
	logger        *logger.BaseLogger
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

	data, err := p.pluginService.InstallLocalPlugin(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
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

func (p *Controller) Activate(ctx fiber.Ctx) error {
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

func (p *Controller) Deactivate(ctx fiber.Ctx) error {
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

func (p *Controller) Status(ctx fiber.Ctx) error {
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

func (p *Controller) Upload(c fiber.Ctx) error {
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

func (p *Controller) InstallFromAttachment(c fiber.Ctx) error {
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
func NewPluginController(pluginService *Service, logger *logger.BaseLogger) *Controller {
	return &Controller{
		pluginService: pluginService,
		logger:        logger.With("module", "plugin-controller"),
	}
}
