package repositories

import (
	"context"

	"github.com/smtdfc/nagare/core/persistence/database"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func (r *SessionRepository) Create(ctx context.Context) (*models.Session, error) {
	session := &models.Session{}
	err := gorm.G[models.Session](r.db).Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func NewSessionRepository() *SessionRepository {
	db, _ := database.GetDatabase()
	return &SessionRepository{
		db: db,
	}
}
