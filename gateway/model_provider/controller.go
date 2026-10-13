package model_provider

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/utils"
)

type ModelProviderController struct {
	modelProviderService *ModelProviderService
}

func (l *ModelProviderController) List(ctx fiber.Ctx) error {
	data, err := l.modelProviderService.ListProviders(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (l *ModelProviderController) Details(ctx fiber.Ctx) error {
	providerID := ctx.Query("provider")
	data, err := l.modelProviderService.GetProviderDetails(ctx, providerID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (l *ModelProviderController) Add(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.AddModelProviderRequest](ctx)
	if err != nil {
		return err
	}

	data, err := l.modelProviderService.AddProvider(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (l *ModelProviderController) Delete(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.DeleteModelProviderRequest](ctx)
	if err != nil {
		return err
	}

	err = l.modelProviderService.DeleteProvider(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, struct{}{}, 200)
}

func (l *ModelProviderController) GetModels(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.GetModelProviderModelsRequest](ctx)
	if err != nil {
		return err
	}

	data, err := l.modelProviderService.GetModels(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

// @Injectable
func NewModelProviderController(
	modelProviderService *ModelProviderService,
) *ModelProviderController {
	return &ModelProviderController{
		modelProviderService: modelProviderService,
	}
}
