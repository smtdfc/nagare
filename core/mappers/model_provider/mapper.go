//go:build !goverter

package model_provider

import (
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/model_provider/provider"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

type ModelProviderMapper struct {
	converter mappers.ModelProviderConverter
}

// @Injectable
func NewModelProviderMapper(credentialMapper *CredentialMapper) *ModelProviderMapper {
	return &ModelProviderMapper{
		converter: &ModelProviderConverterImpl{},
	}
}

func (l *ModelProviderMapper) ToDomain(entity *entities.ModelProvider) *provider.ModelProviderConfig {
	return l.converter.ToDomain(entity)
}

func (l *ModelProviderMapper) ToEntity(domain *provider.ModelProviderConfig) *entities.ModelProvider {
	return l.converter.ToEntity(domain)
}

func (l *ModelProviderMapper) ToProviderInfo(entity *entities.ModelProvider) *provider.ModelProviderInfo {
	return l.converter.ToProviderInfo(entity)
}

func (l *ModelProviderMapper) ToDomains(entities []*entities.ModelProvider) []*provider.ModelProviderConfig {
	return l.converter.ToDomains(entities)
}
