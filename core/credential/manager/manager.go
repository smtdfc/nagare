package manager

import (
	"context"

	"github.com/smtdfc/nagare/core/credential"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
)

type CredentialManager struct {
	credentialRepo   *repositories.CredentialRepository
	credentialMapper *mappers.CredentialMapper
}

func (c *CredentialManager) GetAll(ctx context.Context) ([]*credential.Credential, error) {
	entities, err := c.credentialRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	credentials := make([]*credential.Credential, 0, len(entities))
	for _, entity := range entities {
		credentials = append(credentials, c.credentialMapper.ToDomain(entity))
	}

	return credentials, nil
}

func (c *CredentialManager) GetByID(ctx context.Context, id string) (*credential.Credential, error) {
	entity, err := c.credentialRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return c.credentialMapper.ToDomain(entity), nil
}

func (c *CredentialManager) Create(ctx context.Context, value *credential.Credential) (*credential.Credential, error) {
	entity, err := c.credentialRepo.Add(ctx, c.credentialMapper.ToEntity(value))
	if err != nil {
		return nil, err
	}

	return c.credentialMapper.ToDomain(entity), nil
}

func (c *CredentialManager) Update(ctx context.Context, value *credential.Credential) (*credential.Credential, error) {
	entity := c.credentialMapper.ToEntity(value)
	if err := c.credentialRepo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return c.credentialMapper.ToDomain(entity), nil
}

func (c *CredentialManager) Delete(ctx context.Context, id string) error {
	return c.credentialRepo.DeleteByID(ctx, id)
}

// @Injectable
func NewCredentialManager(credentialRepo *repositories.CredentialRepository, credentialMapper *mappers.CredentialMapper) *CredentialManager {
	return &CredentialManager{
		credentialRepo:   credentialRepo,
		credentialMapper: credentialMapper,
	}
}
