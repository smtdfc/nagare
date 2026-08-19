package repositories

import (
	"context"

	"github.com/smtdfc/nagare/core/persistence"
	"github.com/smtdfc/nagare/core/persistence/database"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func (r *SessionRepository) CreateWithModel(ctx context.Context, session *models.Session) (*models.Session, error) {
	err := gorm.G[models.Session](r.db).Create(ctx, session)
	if err != nil {
		persistence.PersistenceLogger.Error("Failed to create session", "error", err)
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) Create(ctx context.Context) (*models.Session, error) {
	session := &models.Session{}
	err := gorm.G[models.Session](r.db).Create(ctx, session)
	if err != nil {
		persistence.PersistenceLogger.Error("Failed to create session", "error", err)
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) FindByID(id string) (*models.Session, error) {
	var session models.Session
	if err := r.db.Where("id = ?", id).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		persistence.PersistenceLogger.Error("Failed to find session", "error", err)
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) FindByUserID(id string) (*models.Session, error) {
	var session models.Session
	if err := r.db.Where("user_id = ?", id).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		persistence.PersistenceLogger.Error("Failed to find session", "error", err)
		return nil, err
	}
	return &session, nil
}

func NewSessionRepository() *SessionRepository {
	db, _ := database.GetDatabase()
	return &SessionRepository{
		db: db,
	}
}
