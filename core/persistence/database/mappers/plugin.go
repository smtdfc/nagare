package mappers

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"github.com/smtdfc/nagare/shared/plugin"
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
		Features:   plugin.ParseFeaturesFromString(model.Features),
	}
}

func (m *PluginMapper) ToModel(domain *domains.PluginInfo) *models.Plugin {
	var parsedUUID uuid.UUID
	if domain.ID != "" {
		if id, err := uuid.Parse(domain.ID); err == nil {
			parsedUUID = id
		}
	}

	return &models.Plugin{
		ID:         parsedUUID,
		PluginID:   domain.PluginID,
		Name:       domain.Name,
		Active:     domain.Active,
		ApiVersion: domain.ApiVersion,
		Author:     domain.Author,
		Version:    domain.Version,
		Bin:        domain.Bin,
		Features:   domain.Features.ToString(),
	}
}
