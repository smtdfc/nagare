package repositories

import (
	"context"
	"errors"
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

func (p *PluginRepository) FindById(ctx context.Context, id string) (*entities.Plugin, error) {
	var plugin entities.Plugin
	err := p.db.WithContext(ctx).
		Where("id = ?", id).
		First(&plugin).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		p.logger.Error("Failed to get plugin by ID", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get plugin by ID: %w", err)
	}

	return &plugin, nil
}

func (p *PluginRepository) DeleteById(ctx context.Context, id string) error {
	err := p.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", id).
		Delete(&entities.Plugin{}).Error
	if err != nil {
		p.logger.Error("Failed to delete plugin", "error", err, "id", id)
		return fmt.Errorf("failed to delete plugin: %w", err)
	}
	return nil
}

func (p *PluginRepository) FindAll(ctx context.Context) ([]*entities.Plugin, error) {
	var plugins []*entities.Plugin

	err := p.db.WithContext(ctx).
		Find(&plugins).Error

	if err != nil {
		p.logger.Error("Failed to get all plugin", "error", err)
		return nil, fmt.Errorf("failed to get all plugin: %w", err)
	}

	return plugins, nil
}

func (p *PluginRepository) FindActive(ctx context.Context) ([]*entities.Plugin, error) {
	var plugins []*entities.Plugin

	err := p.db.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&plugins).Error

	if err != nil {
		p.logger.Error("Failed to get all active plugin", "error", err)
		return nil, fmt.Errorf("failed to get all active plugin: %w", err)
	}

	return plugins, nil
}

func (p *PluginRepository) FindByPluginId(ctx context.Context, pluginId string) (*entities.Plugin, error) {
	var plugin entities.Plugin

	err := p.db.WithContext(ctx).
		Where("plugin_id = ?", pluginId).
		First(&plugin).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		p.logger.Error("Failed to get  plugin", "error", err, "pluginId", pluginId)
		return nil, fmt.Errorf("failed to get plugin: %w", err)
	}

	return &plugin, nil
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
