package mappers

import (
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/shared/helpers"
)

type PluginMapper struct{}

func (p *PluginMapper) ToDomain(entity *entities.Plugin) *plugin.Plugin {
	return &plugin.Plugin{
		ID:       entity.ID,
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
	return &entities.Plugin{
		ID:       domain.ID,
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
