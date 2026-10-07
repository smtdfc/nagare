package credential

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/utils"
)

type CredentialController struct {
	credentialService *CredentialService
}

func (c *CredentialController) List(ctx fiber.Ctx) error {
	data, err := c.credentialService.ListCredentials(ctx)
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *CredentialController) Details(ctx fiber.Ctx) error {
	data, err := c.credentialService.GetCredentialDetails(ctx, ctx.Query("credential"))
	if err != nil {
		return err
	}

	return utils.ResponseSuccess(ctx, data, 200)
}

func (c *CredentialController) Add(ctx fiber.Ctx) error {
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

func (c *CredentialController) Update(ctx fiber.Ctx) error {
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

func (c *CredentialController) Delete(ctx fiber.Ctx) error {
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
func NewCredentialController(credentialService *CredentialService) *CredentialController {
	return &CredentialController{
		credentialService: credentialService,
	}
}
