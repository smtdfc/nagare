package mappers

import (
	"github.com/google/uuid"

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
		ID:        entity.ID.String(),
		Title:     entity.Title,
		OwnerID:   entity.OwnerID.String(),
		OwnerType: session.GetOwnerType(entity.OwnerType),
		IsArchive: entity.IsArchive,
		ChannelID: entity.ChannelID,
	}
}

func (s *SessionMapper) ToEntity(domain *session.SessionInfo) *entities.Session {
	if domain == nil {
		return nil
	}

	id, _ := uuid.Parse(domain.ID)
	ownerID, _ := uuid.Parse(domain.OwnerID)

	return &entities.Session{
		ID:        id,
		Title:     domain.Title,
		OwnerID:   ownerID,
		OwnerType: string(domain.OwnerType),
		IsArchive: domain.IsArchive,
		ChannelID: domain.ChannelID,
	}
}

func (s *SessionMapper) ToDomains(entities []*entities.Session) []*session.SessionInfo {
	return helpers.SliceMap(entities, s.ToDomain)
}
