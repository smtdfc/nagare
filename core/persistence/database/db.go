package database

import (
	"path"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/persistence"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"github.com/smtdfc/nagare/shared/paths"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase() (*gorm.DB, error) {
	persistence.PersistenceLogger.Info("Initing database ")
	if db != nil {
		return db, nil
	}

	var err error
	db, err = gorm.Open(sqlite.Open(path.Join(paths.DatabaseDir, "nagare.db")), &gorm.Config{})
	if err != nil {
		persistence.PersistenceLogger.Error("Failed to init database", "error", err)
		return nil, err
	}

	err = db.AutoMigrate(&models.Session{}, &models.Message{}, &models.KV{}, &models.LLMProvider{})
	if err != nil {
		persistence.PersistenceLogger.Error("Failed to migrate database", "error", err)
		return nil, err
	}

	return db, nil
}

func GetDatabase() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}

	return nil, custom_errors.NewDatabaseError("The database connection has not been established. Ensure InitDatabase is called before performing operations")
}
