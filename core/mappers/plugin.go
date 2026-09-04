package mappers

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/shared/helpers"
)

type PluginMapper struct{}

func (p *PluginMapper) ToDomain(entity *entities.Plugin) *plugin.Plugin {
	return &plugin.Plugin{
		ID:       entity.ID.String(),
		PluginID: entity.PluginID,
		Name:     entity.Name,
		Author:   entity.Author,
		Features: plugin.ParseFeatureString(entity.Features),
		Version:  entity.Version,
		Bin:      entity.Bin,
		IsActive: entity.IsActive,
	}
}

func (p *PluginMapper) ToEntity(domain *plugin.Plugin) *entities.Plugin {
	var id uuid.UUID
	id, err := uuid.Parse(domain.ID)
	if err != nil {
		id = uuid.Nil
	}

	return &entities.Plugin{
		ID:       id,
		PluginID: domain.PluginID,
		Name:     domain.Name,
		Author:   domain.Author,
		Features: domain.ToFeaturesString(),
		Version:  domain.Version,
		Bin:      domain.Bin,
		IsActive: domain.IsActive,
	}
}

func (p *PluginMapper) ToDomains(entities []*entities.Plugin) []*plugin.Plugin {
	return helpers.SliceMap(entities, p.ToDomain)
}

// @Injectable
func NewPluginMapper() *PluginMapper {
	return &PluginMapper{}
}
