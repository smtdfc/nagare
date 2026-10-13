package database

import (
	"path/filepath"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/pkgs/paths"
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

	err = db.AutoMigrate(
		&entities.KV{},
		&entities.Session{},
		&entities.Message{},
		&entities.Plugin{},
		&entities.ModelProvider{},
		&entities.Task{},
		&entities.Plugin{},
		&entities.Credential{},
	)
	if err != nil {
		databaseLogger.Error("Failed to migrate database", "error", err)
		return nil, err
	}

	migrateModelProviderSchema(db, databaseLogger)

	return db, nil
}

func migrateModelProviderSchema(db *gorm.DB, databaseLogger *logger.BaseLogger) {
	if db.Migrator().HasColumn(&entities.Session{}, "llm_provider_id") {
		if err := db.Exec("UPDATE sessions SET model_provider_id = llm_provider_id WHERE model_provider_id IS NULL OR model_provider_id = ''").Error; err != nil {
			databaseLogger.Error("Failed to migrate session model_provider_id", "error", err)
		}
	}
	if db.Migrator().HasTable("embedding_providers") {
		if err := db.Migrator().DropTable("embedding_providers"); err != nil {
			databaseLogger.Error("Failed to drop legacy embedding_providers table", "error", err)
		}
	}
}
