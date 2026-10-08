package database

import (
	"testing"

	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAutoMigrate_InMemory(t *testing.T) {
	// Verify database entities schema migration works on SQLite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	err = db.AutoMigrate(
		&entities.KV{},
		&entities.Session{},
		&entities.Message{},
		&entities.Plugin{},
		&entities.LLMProvider{},
		&entities.Task{},
		&entities.Credential{},
	)
	if err != nil {
		t.Fatalf("failed to migrate database entities: %v", err)
	}
}
