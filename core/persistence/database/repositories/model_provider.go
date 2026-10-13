package repositories

import (
	"context"
	"errors"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/gorm"
)

type ModelProviderRepository struct {
	db     *gorm.DB
	logger *logger.BaseLogger
}

func (l *ModelProviderRepository) FindAll(ctx context.Context) ([]*entities.ModelProvider, error) {
	var providers []*entities.ModelProvider
	err := l.db.WithContext(ctx).Find(&providers).Error
	if err != nil {
		l.logger.Error("Failed to get all provider", "err", err)
		return nil, err
	}

	return providers, nil
}

func (l *ModelProviderRepository) FindByID(ctx context.Context, id string) (*entities.ModelProvider, error) {
	var provider entities.ModelProvider
	err := l.db.Where("id = ?", id).WithContext(ctx).First(&provider).Preload("Credential").Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		l.logger.Error("Failed to get provider", "id", id, "err", err)
		return nil, err
	}

	return &provider, nil
}

func (l *ModelProviderRepository) DeleteByID(ctx context.Context, id string) error {
	err := l.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ModelProvider{}).Error
	if err != nil {
		l.logger.Error("Failed to delete provider", "id", id, "err", err)
		return err
	}

	return nil
}

func (l *ModelProviderRepository) Update(ctx context.Context, provider *entities.ModelProvider) error {
	result := l.db.WithContext(ctx).Model(provider).Where("id = ?", provider.ID).Updates(provider)

	if result.Error != nil {
		l.logger.Error("Failed to update provider", "id", provider.ID, "err", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		l.logger.Warn("Provider not found for update", "id", provider.ID)
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (l *ModelProviderRepository) Add(ctx context.Context, provider *entities.ModelProvider) (*entities.ModelProvider, error) {
	err := l.db.WithContext(ctx).Create(&provider).Error
	if err != nil {
		l.logger.Error("Failed to create provider", "err", err)
		return nil, err
	}

	return provider, nil
}

// @Injectable
func NewModelProviderRepository(db *gorm.DB, logger *logger.BaseLogger) *ModelProviderRepository {
	return &ModelProviderRepository{
		db:     db,
		logger: logger.With("module", "model-provider-repository"),
	}
}
