package mappers

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/persistence/database/models"
)

type PluginMapper struct{}

func (m *PluginMapper) ToDomain(model *models.Plugin) *domains.PluginInfo {
	return &domains.PluginInfo{
		ID:         model.ID.String(),
		PluginID:   model.PluginID,
		Name:       model.Name,
		Active:     model.Active,
		ApiVersion: model.ApiVersion,
		Author:     model.Author,
		Version:    model.Version,
		Bin:        model.Bin,
	}
}

func (m *PluginMapper) ToModel(domain *domains.PluginInfo) *models.Plugin {
	return &models.Plugin{
		ID:         uuid.MustParse(domain.ID),
		PluginID:   domain.PluginID,
		Name:       domain.Name,
		Active:     domain.Active,
		ApiVersion: domain.ApiVersion,
		Author:     domain.Author,
		Version:    domain.Version,
	}
}
