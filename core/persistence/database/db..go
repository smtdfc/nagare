package database

import (
	"path/filepath"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/shared/paths"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @Injectable
func InitDatabase(logger *logger.BaseLogger) (*gorm.DB, error) {
	var databaseLogger = logger.With("module", "database")
	var err error
	db, err := gorm.Open(sqlite.Open(filepath.Join(paths.DatabaseDir, "nagare.db")), &gorm.Config{})
	if err != nil {
		databaseLogger.Error("Failed to init database", "error", err)
		return nil, err
	}

	err = db.AutoMigrate(&entities.KV{}, &entities.Session{}, &entities.Message{}, &entities.Plugin{}, &entities.LLMProvider{})
	if err != nil {
		databaseLogger.Error("Failed to migrate database", "error", err)
		return nil, err
	}

	return db, nil
}
