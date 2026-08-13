package repositories

import (
	"github.com/smtdfc/nagare/core/persistence/database"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KVRepository struct {
	db *gorm.DB
}

func (r *KVRepository) GetAllKeyByTarget(target string) ([]models.KV, error) {
	var kv []models.KV
	if err := r.db.Where("target = ?", target).Find(&kv).Error; err != nil {
		return nil, err
	}
	return kv, nil
}

func (r *KVRepository) Save(kv []models.KV) error {
	if len(kv) == 0 {
		return nil
	}

	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&kv).Error

	return err
}

func NewKVRepository() *KVRepository {
	db, _ := database.GetDatabase()
	return &KVRepository{db: db}
}
