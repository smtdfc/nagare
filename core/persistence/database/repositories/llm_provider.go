package repositories

import (
	"github.com/smtdfc/nagare/core/persistence"
	"github.com/smtdfc/nagare/core/persistence/database"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"gorm.io/gorm"
)

type LLMProviderRepository struct {
	db *gorm.DB
}

func (r *LLMProviderRepository) FindAll() ([]models.LLMProvider, error) {
	var providers []models.LLMProvider
	if err := r.db.Find(&providers).Error; err != nil {
		persistence.PersistenceLogger.Error("Failed to find all providers", "error", err)
		return nil, err
	}
	return providers, nil
}

func (r *LLMProviderRepository) Save(provider *models.LLMProvider) error {
	if err := r.db.Save(provider).Error; err != nil {
		persistence.PersistenceLogger.Error("Failed to save provider", "error", err)
		return err
	}
	return nil
}

func (r *LLMProviderRepository) UpdateByID(id string, provider *models.LLMProvider) error {
	if err := r.db.Model(&models.LLMProvider{}).Where("id = ?", id).Updates(provider).Error; err != nil {
		persistence.PersistenceLogger.Error("Failed to update provider", "error", err)
		return err
	}
	return nil
}

func (r *LLMProviderRepository) FindByID(id string) (*models.LLMProvider, error) {
	var provider models.LLMProvider
	if err := r.db.Where("id = ?", id).First(&provider).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		persistence.PersistenceLogger.Error("Failed to find provider", "error", err)
		return nil, err
	}
	return &provider, nil
}

func (r *LLMProviderRepository) CreateProvider(provider *models.LLMProvider) error {
	if err := r.db.Create(provider).Error; err != nil {
		persistence.PersistenceLogger.Error("Failed to create provider", "error", err)
		return err
	}
	return nil
}

func (r *LLMProviderRepository) DeleteByID(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&models.LLMProvider{}).Error; err != nil {
		persistence.PersistenceLogger.Error("Failed to delete provider", "error", err)
		return err
	}
	return nil
}

func NewLLMProviderRepository() *LLMProviderRepository {
	db, _ := database.GetDatabase()
	return &LLMProviderRepository{db: db}
}
