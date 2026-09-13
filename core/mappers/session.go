package mappers

import (
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/shared/helpers"
)

type SessionMapper struct{}

// @Injectable
func NewSessionMapper() *SessionMapper {
	return &SessionMapper{}
}

func (s *SessionMapper) ToDomain(entity *entities.Session) *session.SessionInfo {
	if entity == nil {
		return nil
	}
	return &session.SessionInfo{
		ID:        entity.ID,
		Title:     entity.Title,
		OwnerID:   entity.OwnerID,
		OwnerType: session.GetOwnerType(entity.OwnerType),
		IsArchive: entity.IsArchive,
		ChannelID: entity.ChannelID,
	}
}

func (s *SessionMapper) ToEntity(domain *session.SessionInfo) *entities.Session {
	return &entities.Session{
		ID:        domain.ID,
		Title:     domain.Title,
		OwnerID:   domain.OwnerID,
		OwnerType: string(domain.OwnerType),
		IsArchive: domain.IsArchive,
		ChannelID: domain.ChannelID,
	}
}

func (s *SessionMapper) ToDomains(entities []*entities.Session) []*session.SessionInfo {
	return helpers.SliceMap(entities, s.ToDomain)
}
