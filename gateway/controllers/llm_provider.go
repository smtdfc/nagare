package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/gateway/services"
	"github.com/smtdfc/nagare/gateway/utils"
	"github.com/smtdfc/nagare/shared/dtos/rest"
)

type LLMProviderController struct {
	llmProviderService *services.LLMProviderService
}

func (l *LLMProviderController) List(ctx fiber.Ctx) error {
	data, err := l.llmProviderService.ListProviders(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (l *LLMProviderController) Details(ctx fiber.Ctx) error {
	providerID := ctx.Query("provider")
	data, err := l.llmProviderService.GetProviderDetails(ctx, providerID)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (l *LLMProviderController) Add(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.AddLLMProviderRequest](ctx)
	if err != nil {
		return err
	}

	data, err := l.llmProviderService.AddProvider(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (l *LLMProviderController) Delete(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.DeleteLLMProviderRequest](ctx)
	if err != nil {
		return err
	}

	err = l.llmProviderService.DeleteProvider(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, struct{}{}, 200)
}

func (l *LLMProviderController) GetModels(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.GetLLMProviderModelsRequest](ctx)
	if err != nil {
		return err
	}

	data, err := l.llmProviderService.GetModels(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

// @Injectable
func NewLLMProviderController(
	llmProviderService *services.LLMProviderService,
) *LLMProviderController {
	return &LLMProviderController{
		llmProviderService: llmProviderService,
	}
}
