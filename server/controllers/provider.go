package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/server/services"
	"github.com/smtdfc/nagare/server/utils"
	"github.com/smtdfc/nagare/shared/dto"
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

func (p *ProviderController) CreateProvider(ctx fiber.Ctx) error {
	var data dto.CreateProviderRequest
	if err := ctx.Bind().JSON(&data); err != nil {
		return err
	}

	fmt.Println(data)
	resp, err := p.providerService.CreateProvider(data)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(resp, 200, ctx)
}

func (p *ProviderController) UpdateProvider(ctx fiber.Ctx) error {
	var data dto.UpdateProviderRequest
	if err := ctx.Bind().JSON(&data); err != nil {
		return err
	}

	resp, err := p.providerService.UpdateProvider(data)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(resp, 200, ctx)
}

func (p *ProviderController) DeleteProvider(ctx fiber.Ctx) error {
	var data dto.DeleteProviderRequest
	if err := ctx.Bind().JSON(&data); err != nil {
		return err
	}

	resp, err := p.providerService.DeleteProvider(data)
	if err != nil {
		return err
	}

	return utils.SuccessResponse(resp, 200, ctx)
}

func (p *ProviderController) FetchModel(ctx fiber.Ctx) error {
	var data dto.FetchModelRequest
	if err := ctx.Bind().JSON(&data); err != nil {
		return err
	}

	resp, err := p.providerService.FetchModel(data)
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
