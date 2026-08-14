package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/services"
	"github.com/smtdfc/nagare/server/utils"
)

type PluginController struct {
	pluginService *services.PluginService
}

func (c *PluginController) GetAll(ctx fiber.Ctx) error {
	data, err := c.pluginService.GetAllPlugin()
	if err != nil {
		return err
	}
	return utils.SuccessResponse(data, 200, ctx)
}

// @Injectable
func NewPluginController(pluginService *services.PluginService) *PluginController {
	return &PluginController{pluginService: pluginService}
}
