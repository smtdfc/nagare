package repositories

import (
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
		return nil, err
	}
	return providers, nil
}

func (r *LLMProviderRepository) Save(provider *models.LLMProvider) error {
	if err := r.db.Save(provider).Error; err != nil {
		return err
	}
	return nil
}

func (r *LLMProviderRepository) FindByID(id string) (*models.LLMProvider, error) {
	var provider models.LLMProvider
	if err := r.db.First(&provider, id).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

func NewLLMProviderRepository() *LLMProviderRepository {
	db, _ := database.GetDatabase()
	return &LLMProviderRepository{db: db}
}
