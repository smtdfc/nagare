package credential

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/utils"
)

type Controller struct {
	credentialService *Service
}

func (c *Controller) List(ctx fiber.Ctx) error {
	data, err := c.credentialService.ListCredentials(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *Controller) Details(ctx fiber.Ctx) error {
	data, err := c.credentialService.GetCredentialDetails(ctx, ctx.Query("credential"))
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *Controller) Add(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.AddCredentialRequest](ctx)
	if err != nil {
		return err
	}

	data, err := c.credentialService.AddCredential(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *Controller) Update(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.UpdateCredentialRequest](ctx)
	if err != nil {
		return err
	}

	data, err := c.credentialService.UpdateCredential(ctx, request)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *Controller) Delete(ctx fiber.Ctx) error {
	request, err := utils.ParseBody[*rest.DeleteCredentialRequest](ctx)
	if err != nil {
		return err
	}

	if err := c.credentialService.DeleteCredential(ctx, request); err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, struct{}{}, 200)
}

// @Injectable
func NewCredentialController(credentialService *Service) *Controller {
	return &Controller{
		credentialService: credentialService,
	}
}
