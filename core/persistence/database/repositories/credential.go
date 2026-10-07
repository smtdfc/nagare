package repositories

import (
	"context"
	"errors"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/gorm"
)

type CredentialRepository struct {
	db     *gorm.DB
	logger *logger.BaseLogger
}

func (c *CredentialRepository) FindAll(ctx context.Context) ([]*entities.Credential, error) {
	var credentials []*entities.Credential
	if err := c.db.WithContext(ctx).Find(&credentials).Error; err != nil {
		c.logger.Error("Failed to get all credentials", "error", err)
		return nil, err
	}

	return credentials, nil
}

func (c *CredentialRepository) FindByID(ctx context.Context, id string) (*entities.Credential, error) {
	var credential entities.Credential
	if err := c.db.WithContext(ctx).Where("id = ?", id).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		c.logger.Error("Failed to get credential", "id", id, "error", err)
		return nil, err
	}

	return &credential, nil
}

func (c *CredentialRepository) Add(ctx context.Context, credential *entities.Credential) (*entities.Credential, error) {
	if err := c.db.WithContext(ctx).Create(credential).Error; err != nil {
		c.logger.Error("Failed to create credential", "error", err)
		return nil, err
	}

	return credential, nil
}

func (c *CredentialRepository) Update(ctx context.Context, credential *entities.Credential) error {
	result := c.db.WithContext(ctx).Model(credential).Where("id = ?", credential.ID).Updates(credential)
	if result.Error != nil {
		c.logger.Error("Failed to update credential", "id", credential.ID, "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (c *CredentialRepository) DeleteByID(ctx context.Context, id string) error {
	result := c.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Credential{})
	if result.Error != nil {
		c.logger.Error("Failed to delete credential", "id", id, "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// @Injectable
func NewCredentialRepository(db *gorm.DB, logger *logger.BaseLogger) *CredentialRepository {
	return &CredentialRepository{
		db:     db,
		logger: logger.With("module", "credential-repository"),
	}
}
