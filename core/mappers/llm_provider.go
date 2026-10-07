package mappers

import (
	"strings"

	"github.com/smtdfc/nagare/core/credential"
	"github.com/smtdfc/nagare/core/llm/provider"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

type LLMProviderMapper struct {
	credentialMapper *CredentialMapper
}

func (l *LLMProviderMapper) ToDomain(entity *entities.LLMProvider) *provider.LLMProviderConfig {
	var credentialValue *credential.Credential
	if entity == nil {
		return nil
	}

	if entity.Credential != nil {
		credentialValue = l.credentialMapper.ToDomain(entity.Credential)
	}

	return &provider.LLMProviderConfig{
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

func (l *LLMProviderMapper) ToEntity(domain *provider.LLMProviderConfig) *entities.LLMProvider {
	return &entities.LLMProvider{
		ID:           domain.ID,
		Name:         domain.Name,
		Compatible:   domain.Compatible.ToString(),
		ApiKey:       domain.ApiKey,
		Models:       strings.Join(domain.Models, ","),
		BaseURL:      domain.BaseURL,
		CredentialID: domain.CredentialID,
	}
}

func (l *LLMProviderMapper) ToProviderInfo(entity *entities.LLMProvider) *provider.LLMProviderInfo {
	return &provider.LLMProviderInfo{
		ID:         entity.ID,
		Name:       entity.Name,
		Compatible: provider.GetCompatibleFromString(entity.Compatible),
	}
}

func (l *LLMProviderMapper) ToDomains(entities []*entities.LLMProvider) []*provider.LLMProviderConfig {
	return helpers.SliceMap(entities, l.ToDomain)
}

// @Injectable
func NewLLMProviderMapper(credentialMapper *CredentialMapper) *LLMProviderMapper {
	return &LLMProviderMapper{
		credentialMapper: credentialMapper,
	}
}
