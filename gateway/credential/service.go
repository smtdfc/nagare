package credential

import (
	"context"

	"github.com/google/uuid"
	core_credential "github.com/smtdfc/nagare/core/credential"
	credential_manager "github.com/smtdfc/nagare/core/credential/manager"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

func toCredentialDTO(value *core_credential.Credential) *rest.Credential {
	if value == nil {
		return nil
	}

	return &rest.Credential{
		ID:   value.ID.String(),
		Name: value.Name,
	}
}

type Service struct {
	credentialManager *credential_manager.CredentialManager
}

func (c *Service) ListCredentials(ctx context.Context) (*rest.GetListCredentialResponse, error) {
	credentials, err := c.credentialManager.GetAll(ctx)
	if err != nil {
		return nil, custom_errors.ErrGetAllCredentialFailed
	}

	return &rest.GetListCredentialResponse{
		Credentials: helpers.SliceMap(credentials, toCredentialDTO),
	}, nil
}

func (c *Service) GetCredentialDetails(ctx context.Context, id string) (*rest.GetCredentialDetailsResponse, error) {
	credentialValue, err := c.credentialManager.GetByID(ctx, id)
	if err != nil {
		return nil, custom_errors.ErrGetCredentialFailed
	}
	if credentialValue == nil {
		return nil, custom_errors.ErrCredentialNotFound
	}

	return &rest.GetCredentialDetailsResponse{
		Credential: toCredentialDTO(credentialValue),
	}, nil
}

func (c *Service) AddCredential(ctx context.Context, request *rest.AddCredentialRequest) (*rest.AddCredentialResponse, error) {
	credentialValue, err := c.credentialManager.Create(ctx, &core_credential.Credential{
		Name:   request.Name,
		ApiKey: request.ApiKey,
	})
	if err != nil {
		return nil, custom_errors.ErrCreateCredentialFailed
	}

	return &rest.AddCredentialResponse{
		Credential: toCredentialDTO(credentialValue),
	}, nil
}

func (c *Service) UpdateCredential(ctx context.Context, request *rest.UpdateCredentialRequest) (*rest.UpdateCredentialResponse, error) {
	id, err := uuid.Parse(request.ID)
	if err != nil {
		return nil, custom_errors.ErrUpdateCredentialFailed
	}

	credentialValue, err := c.credentialManager.Update(ctx, &core_credential.Credential{
		ID:     id,
		Name:   request.Name,
		ApiKey: request.ApiKey,
	})
	if err != nil {
		return nil, custom_errors.ErrUpdateCredentialFailed
	}

	return &rest.UpdateCredentialResponse{
		Credential: toCredentialDTO(credentialValue),
	}, nil
}

func (c *Service) DeleteCredential(ctx context.Context, request *rest.DeleteCredentialRequest) error {
	if err := c.credentialManager.Delete(ctx, request.ID); err != nil {
		return custom_errors.ErrDeleteCredentialFailed
	}

	return nil
}

// @Injectable
func NewCredentialService(credentialManager *credential_manager.CredentialManager) *Service {
	return &Service{
		credentialManager: credentialManager,
	}
}
