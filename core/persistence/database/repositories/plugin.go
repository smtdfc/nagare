package repositories

import (
	"errors"
	"fmt"

	"github.com/smtdfc/nagare/core/persistence"
	"github.com/smtdfc/nagare/core/persistence/database"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"gorm.io/gorm"
)

type PluginRepository struct {
	db *gorm.DB
}

func (r *PluginRepository) CreatePlugin(plugin *models.Plugin) error {
	result := r.db.Create(plugin)
	if result.Error != nil {
		persistence.PersistenceLogger.Error("failed to create plugin", "error", result.Error)
		return result.Error
	}
	return nil
}

func (r *PluginRepository) GetAllActivePlugins() ([]models.Plugin, error) {
	var plugins []models.Plugin
	result := r.db.Where("active = ?", true).Find(&plugins)
	if result.Error != nil {
		persistence.PersistenceLogger.Error("failed to get active plugins", "error", result.Error)
		return nil, result.Error
	}
	return plugins, nil
}

func (r *PluginRepository) GetAllPlugins() ([]models.Plugin, error) {
	var plugins []models.Plugin
	result := r.db.Find(&plugins)
	if result.Error != nil {
		persistence.PersistenceLogger.Error("failed to get all plugins", "error", result.Error)
		return nil, result.Error
	}
	return plugins, nil
}

func (r *PluginRepository) DeletePluginByID(id string) error {
	result := r.db.Delete(&models.Plugin{}, id)
	if result.Error != nil {
		persistence.PersistenceLogger.Error("failed to delete plugin", "error", result.Error)
		return result.Error
	}
	return nil
}

func (r *PluginRepository) GetPluginByID(id string) (*models.Plugin, error) {
	var plugin models.Plugin
	result := r.db.First(&plugin, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			persistence.PersistenceLogger.Warn("plugin not found", "id", id)
			return nil, fmt.Errorf("plugin with id %s not found: %w", id, result.Error)
		}
		persistence.PersistenceLogger.Error("failed to get plugin by id", "error", result.Error, "id", id)
		return nil, result.Error
	}
	return &plugin, nil
}

func NewPluginRepository() *PluginRepository {
	db, _ := database.GetDatabase()
	return &PluginRepository{db: db}
}
