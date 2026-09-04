package repositories

import (
	"context"
	"errors"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KVRepository struct {
	db     *gorm.DB
	logger *logger.BaseLogger
}

// @Injectable
func NewKVRepository(db *gorm.DB, logger *logger.BaseLogger) *KVRepository {
	return &KVRepository{
		db:     db,
		logger: logger.With("module", "kv-repository"),
	}
}

func (r *KVRepository) Upsert(ctx context.Context, kvs []*entities.KV) error {
	if len(kvs) == 0 {
		return nil
	}

	for _, kv := range kvs {
		if kv.Key == "" || kv.Scope == "" {
			return errors.New("key and scope cannot be empty for any item")
		}
	}

	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}, {Name: "scope"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).WithContext(ctx).Create(&kvs).Error

	return err
}

func (r *KVRepository) GetByScope(ctx context.Context, scope string) ([]*entities.KV, error) {
	if scope == "" {
		return nil, errors.New("scope cannot be empty")
	}

	var kvs []*entities.KV
	err := r.db.Where("scope = ?", scope).WithContext(ctx).Find(&kvs).Error
	if err != nil {
		return nil, err
	}

	return kvs, nil
}
