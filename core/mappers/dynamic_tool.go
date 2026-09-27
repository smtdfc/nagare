package mappers

import (
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/shared/helpers"
)

type DynamicToolMapper struct {
	pluginMapper *PluginMapper
}

func (t *DynamicToolMapper) ToDomain(entity *entities.DynamicTool) *tool.DynamicTool {
	return &tool.DynamicTool{
		Name:        entity.Name,
		Description: entity.Description,
		Args:        entity.ArgsSchema,
		PluginID:    entity.PluginID,
		Plugin:      t.pluginMapper.ToDomain(entity.Plugin),
	}
}

func (t *DynamicToolMapper) ToDomains(tools []*entities.DynamicTool) []*tool.DynamicTool {
	return helpers.SliceMap(tools, t.ToDomain)
}

func (t *DynamicToolMapper) ToEntity(domain *tool.DynamicTool) *entities.DynamicTool {
	return &entities.DynamicTool{
		Name:        domain.Name,
		Description: domain.Description,
		ArgsSchema:  domain.Args,
		PluginID:    domain.PluginID,
	}
}

func (t *DynamicToolMapper) ToEntities(domains []*tool.DynamicTool) []*entities.DynamicTool {
	return helpers.SliceMap(domains, func(domain *tool.DynamicTool) *entities.DynamicTool {
		return t.ToEntity(domain)
	})
}

// @Injectable
func NewDynamicToolMapper(pluginMapper *PluginMapper) *DynamicToolMapper {
	return &DynamicToolMapper{
		pluginMapper: pluginMapper,
	}
}
