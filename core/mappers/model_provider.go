package mappers

import (
	"strings"

	"github.com/smtdfc/nagare/core/credential"
	"github.com/smtdfc/nagare/core/model_provider/provider"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

type ModelProviderMapper struct {
	credentialMapper *CredentialMapper
}

func (l *ModelProviderMapper) ToDomain(entity *entities.ModelProvider) *provider.ModelProviderConfig {
	var credentialValue *credential.Credential
	if entity == nil {
		return nil
	}

	if entity.Credential != nil {
		credentialValue = l.credentialMapper.ToDomain(entity.Credential)
	}

	return &provider.ModelProviderConfig{
		ID:           entity.ID,
		Name:         entity.Name,
		Compatible:   provider.GetCompatibleFromString(entity.Compatible),
		ApiKey:       entity.ApiKey,
		Models:       strings.Split(entity.Models, ","),
		BaseURL:      entity.BaseURL,
		CredentialID: entity.CredentialID,
		Credential:   credentialValue,
	}
}

func (l *ModelProviderMapper) ToEntity(domain *provider.ModelProviderConfig) *entities.ModelProvider {
	return &entities.ModelProvider{
		ID:           domain.ID,
		Name:         domain.Name,
		Compatible:   domain.Compatible.ToString(),
		ApiKey:       domain.ApiKey,
		Models:       strings.Join(domain.Models, ","),
		BaseURL:      domain.BaseURL,
		CredentialID: domain.CredentialID,
	}
}

func (l *ModelProviderMapper) ToProviderInfo(entity *entities.ModelProvider) *provider.ModelProviderInfo {
	return &provider.ModelProviderInfo{
		ID:         entity.ID,
		Name:       entity.Name,
		Compatible: provider.GetCompatibleFromString(entity.Compatible),
	}
}

func (l *ModelProviderMapper) ToDomains(entities []*entities.ModelProvider) []*provider.ModelProviderConfig {
	return helpers.SliceMap(entities, l.ToDomain)
}

// @Injectable
func NewModelProviderMapper(credentialMapper *CredentialMapper) *ModelProviderMapper {
	return &ModelProviderMapper{
		credentialMapper: credentialMapper,
	}
}
