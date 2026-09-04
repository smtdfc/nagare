package repositories

import (
	"context"
	"fmt"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PluginRepository struct {
	db     *gorm.DB
	logger *logger.BaseLogger
}

func (p *PluginRepository) GetAllPlugin(ctx context.Context) ([]*entities.Plugin, error) {
	var plugins []*entities.Plugin

	err := p.db.WithContext(ctx).
		Find(&plugins).Error

	if err != nil {
		p.logger.Error("Failed to get all plugin", "error", err)
		return nil, fmt.Errorf("Failed to get all plugin: %w", err)
	}

	return plugins, nil
}

func (p *PluginRepository) GetAllActivePlugin(ctx context.Context) ([]*entities.Plugin, error) {
	var plugins []*entities.Plugin

	err := p.db.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&plugins).Error

	if err != nil {
		p.logger.Error("Failed to get all active plugin", "error", err)
		return nil, fmt.Errorf("Failed to get all active plugin: %w", err)
	}

	return plugins, nil
}

func (p *PluginRepository) CreateOrUpdate(ctx context.Context, plugin *entities.Plugin) (*entities.Plugin, error) {
	err := p.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "plugin_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"author",
				"features",
				"version",
				"bin",
				"is_active",
				"updated_at",
			}),
		}).
		Create(plugin).Error

	if err != nil {
		p.logger.Error("Failed to create or update plugin", "plugin_id", plugin.PluginID, "error", err)
		return nil, fmt.Errorf("failed to create or update plugin: %w", err)
	}

	return plugin, nil
}

// @Injectable
func NewPluginRepository(db *gorm.DB, logger *logger.BaseLogger) *PluginRepository {
	return &PluginRepository{
		db:     db,
		logger: logger.With("module", "plugin-repository"),
	}
}
