package mappers

import (
	"strings"

	"github.com/google/uuid"

	"github.com/smtdfc/nagare/core/llm_provider"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/shared/helpers"
)

type LLMProviderMapper struct {
}

func (l *LLMProviderMapper) ToDomain(entity *entities.LLMProvider) *llm_provider.LLMProviderConfig {
	return &llm_provider.LLMProviderConfig{
		ID:         entity.ID.String(),
		Name:       entity.Name,
		Compatible: llm_provider.GetCompatibleFromString(entity.Compatible),
		ApiKey:     entity.ApiKey,
		Models:     strings.Split(entity.Models, ","),
		BaseURL:    entity.BaseURL,
	}
}

func (l *LLMProviderMapper) ToEntity(domain *llm_provider.LLMProviderConfig) *entities.LLMProvider {
	var providerID uuid.UUID
	if parsedID, err := uuid.Parse(domain.ID); err == nil {
		providerID = parsedID
	}

	return &entities.LLMProvider{
		ID:         providerID,
		Name:       domain.Name,
		Compatible: domain.Compatible.ToString(),
		ApiKey:     domain.ApiKey,
		Models:     strings.Join(domain.Models, ","),
		BaseURL:    domain.BaseURL,
	}
}

func (l *LLMProviderMapper) ToDomains(entities []*entities.LLMProvider) []*llm_provider.LLMProviderConfig {
	return helpers.SliceMap(entities, l.ToDomain)
}

// @Injectable
func NewLLMProviderMapper() *LLMProviderMapper {
	return &LLMProviderMapper{}
}
