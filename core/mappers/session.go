package mappers

import (
	"github.com/smtdfc/nagare/core/llm_provider"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

type SessionMapper struct {
	llmProviderMapper *LLMProviderMapper
}

// @Injectable
func NewSessionMapper(llmProviderMapper *LLMProviderMapper) *SessionMapper {
	return &SessionMapper{
		llmProviderMapper: llmProviderMapper,
	}
}

func (s *SessionMapper) ToDomain(entity *entities.Session) *session.SessionInfo {
	if entity == nil {
		return nil
	}
	var providerInfo *llm_provider.LLMProviderInfo
	if entity.LLMProvider != nil {
		providerInfo = s.llmProviderMapper.ToProviderInfo(entity.LLMProvider)
	}

	return &session.SessionInfo{
		ID:              entity.ID,
		Title:           entity.Title,
		OwnerID:         entity.OwnerID,
		OwnerType:       session.GetOwnerType(entity.OwnerType),
		IsArchive:       entity.IsArchive,
		ChannelID:       entity.ChannelID,
		CurrentLLMModel: entity.CurrentModel,
		LLMProviderID:   entity.LLMProviderID,
		LLMProvider:     providerInfo,
	}
}

func (s *SessionMapper) ToEntity(domain *session.SessionInfo) *entities.Session {
	return &entities.Session{
		ID:            domain.ID,
		Title:         domain.Title,
		OwnerID:       domain.OwnerID,
		OwnerType:     string(domain.OwnerType),
		IsArchive:     domain.IsArchive,
		ChannelID:     domain.ChannelID,
		CurrentModel:  domain.CurrentLLMModel,
		LLMProviderID: domain.LLMProviderID,
	}
}

func (s *SessionMapper) ToDomains(entities []*entities.Session) []*session.SessionInfo {
	return helpers.SliceMap(entities, s.ToDomain)
}
