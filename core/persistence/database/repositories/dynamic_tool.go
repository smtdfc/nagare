package repositories

import (
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DynamicToolRepository struct {
	db     *gorm.DB
	logger *logger.BaseLogger
}

func (d *DynamicToolRepository) Upsert(entity *entities.DynamicTool) error {
	return d.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plugin_id"}, {Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"description", "args_schema", "updated_at"}),
	}).Create(&entity).Error
}

func (d *DynamicToolRepository) UpsertBatch(entities []*entities.DynamicTool) error {
	if len(entities) == 0 {
		return nil
	}

	err := d.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plugin_id"}, {Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"description", "args_schema", "updated_at"}),
	}).Create(&entities).Error
	if err != nil {
		d.logger.Error("Failed to upsert", "err", err)
	}

	return err
}

func (d *DynamicToolRepository) FindByKeywords(keywords []string) ([]*entities.DynamicTool, error) {
	var tools []*entities.DynamicTool

	query := d.db.Model(&entities.DynamicTool{})

	for _, kw := range keywords {
		if kw == "" {
			continue
		}
		searchTerm := "%" + kw + "%"
		query = query.Or("description LIKE ? OR name LIKE ?", searchTerm, searchTerm)
	}

	err := query.Find(&tools).Error
	return tools, err
}

// @Injectable
func NewDynamicToolRepository(
	db *gorm.DB,
	logger *logger.BaseLogger,
) *DynamicToolRepository {
	return &DynamicToolRepository{
		db:     db,
		logger: logger.With("module", "dynamic-tool-repository"),
	}
}
