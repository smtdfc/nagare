package mappers

import (
	"github.com/smtdfc/nagare/core/credential"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

type CredentialMapper struct {
}

// @Injectable
func NewCredentialMapper() *CredentialMapper {
	return &CredentialMapper{}
}

func (m *CredentialMapper) ToDomain(entity *entities.Credential) *credential.Credential {
	if entity == nil {
		return nil
	}

	return &credential.Credential{
		ID:     entity.ID,
		Name:   entity.Name,
		ApiKey: entity.ApiKey,
	}
}

func (m *CredentialMapper) ToEntity(domain *credential.Credential) *entities.Credential {
	if domain == nil {
		return nil
	}

	return &entities.Credential{
		ID:     domain.ID,
		Name:   domain.Name,
		ApiKey: domain.ApiKey,
	}
}
