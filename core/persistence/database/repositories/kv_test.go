package repositories

import (
	"context"
	"log/slog"
	"testing"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, *logger.BaseLogger) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory: %v", err)
	}
	err = db.AutoMigrate(
		&entities.KV{},
		&entities.Credential{},
		&entities.Session{},
		&entities.Message{},
		&entities.Plugin{},
		&entities.LLMProvider{},
		&entities.Task{},
	)
	if err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}
	l := &logger.BaseLogger{Logger: *slog.Default()}
	return db, l
}

func TestKVRepository_UpsertAndFind(t *testing.T) {
	db, l := setupTestDB(t)
	repo := NewKVRepository(db, l)
	ctx := context.Background()

	// Upsert entries
	entriesList := []*entities.KV{
		{Key: "k1", Value: "v1", Scope: "app"},
		{Key: "k2", Value: "v2", Scope: "app"},
	}
	err := repo.Upsert(ctx, entriesList)
	if err != nil {
		t.Fatalf("unexpected error on Upsert: %v", err)
	}

	// Find by scope
	results, err := repo.FindByScope(ctx, "app")
	if err != nil {
		t.Fatalf("unexpected error on FindByScope: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 items, got %d", len(results))
	}

	// Error validation on empty scope
	_, err = repo.FindByScope(ctx, "")
	if err == nil {
		t.Errorf("expected error for empty scope")
	}
}
