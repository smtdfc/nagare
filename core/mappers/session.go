package mappers

import (
	"github.com/smtdfc/nagare/core/model/provider"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

type SessionMapper struct {
	modelProviderMapper *ModelProviderMapper
}

// @Injectable
func NewSessionMapper(modelProviderMapper *ModelProviderMapper) *SessionMapper {
	return &SessionMapper{
		modelProviderMapper: modelProviderMapper,
	}
}

func (s *SessionMapper) ToDomain(entity *entities.Session) *session.Info {
	if entity == nil {
		return nil
	}
	var providerInfo *provider.ModelProviderInfo
	if entity.ModelProvider != nil {
		providerInfo = s.modelProviderMapper.ToProviderInfo(entity.ModelProvider)
	}

	return &session.Info{
		ID:              entity.ID,
		Title:           entity.Title,
		OwnerID:         entity.OwnerID,
		OwnerType:       session.GetOwnerType(entity.OwnerType),
		IsArchive:       entity.IsArchive,
		ChannelID:       entity.ChannelID,
		CurrentModel: entity.CurrentModel,
		ModelProviderID:   entity.ModelProviderID,
		ModelProvider:     providerInfo,
	}
}

func (s *SessionMapper) ToEntity(domain *session.Info) *entities.Session {
	return &entities.Session{
		ID:            domain.ID,
		Title:         domain.Title,
		OwnerID:       domain.OwnerID,
		OwnerType:     string(domain.OwnerType),
		IsArchive:     domain.IsArchive,
		ChannelID:     domain.ChannelID,
		CurrentModel:  domain.CurrentModel,
		ModelProviderID: domain.ModelProviderID,
	}
}

func (s *SessionMapper) ToDomains(entities []*entities.Session) []*session.Info {
	return helpers.SliceMap(entities, s.ToDomain)
}
