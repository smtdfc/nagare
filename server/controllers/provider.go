package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/services"
	"github.com/smtdfc/nagare/server/utils"
)

type ProviderController struct {
	providerService *services.ProviderService
}

func (p *ProviderController) GetListProvider(ctx fiber.Ctx) error {
	resp, err := p.providerService.GetListProvider()
	if err != nil {
		return err
	}

	return utils.SuccessResponse(resp, 200, ctx)
}

func (p *ProviderController) GetProviderDetails(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	resp, err := p.providerService.GetProviderDetails(id)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(resp, 200, ctx)
}

// @Injectable
func NewProviderController(ProviderService *services.ProviderService) *ProviderController {
	return &ProviderController{
		providerService: ProviderService,
	}
}
