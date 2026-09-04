package repositories

import (
	"context"
	"errors"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/gorm"
)

type LLMProviderRepository struct {
	db     *gorm.DB
	logger *logger.BaseLogger
}

func (l *LLMProviderRepository) GetAllProvider(ctx context.Context) ([]*entities.LLMProvider, error) {
	var providers []*entities.LLMProvider
	err := l.db.WithContext(ctx).Find(&providers).Error
	if err != nil {
		l.logger.Error("Failed to get all provider", "err", err)
		return nil, err
	}

	return providers, nil
}

func (l *LLMProviderRepository) GetProviderByID(ctx context.Context, id string) (*entities.LLMProvider, error) {
	var provider entities.LLMProvider
	err := l.db.Where("id = ?", id).WithContext(ctx).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		l.logger.Error("Failed to get provider", "id", id, "err", err)
		return nil, err
	}

	return &provider, nil
}

func (l *LLMProviderRepository) DeleteProviderByID(ctx context.Context, id string) error {
	err := l.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.LLMProvider{}).Error
	if err != nil {
		l.logger.Error("Failed to delete provider", "id", id, "err", err)
		return err
	}

	return nil
}

func (l *LLMProviderRepository) UpdateProvider(ctx context.Context, provider *entities.LLMProvider) error {
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

func (l *LLMProviderRepository) AddProvider(ctx context.Context, provider *entities.LLMProvider) (*entities.LLMProvider, error) {
	err := l.db.WithContext(ctx).Create(&provider).Error
	if err != nil {
		l.logger.Error("Failed to create provider", "err", err)
		return nil, err
	}

	return provider, nil
}

// @Injectable
func NewLLMProviderRepository(db *gorm.DB, logger *logger.BaseLogger) *LLMProviderRepository {
	return &LLMProviderRepository{
		db:     db,
		logger: logger.With("module", "llm-provider-repository"),
	}
}
