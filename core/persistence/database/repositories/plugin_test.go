package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func TestPluginRepository_CRUD(t *testing.T) {
	db, l := setupTestDB(t)
	repo := &PluginRepository{db: db, logger: l}
	ctx := context.Background()

	id := uuid.New()
	p := &entities.Plugin{
		ID:          id,
		PackageName: "com.test.plugin",
		Name:        "Test Plugin",
		Version:     "1.0.0",
		IsActive:    true,
	}

	if err := db.Create(p).Error; err != nil {
		t.Fatalf("failed to insert plugin: %v", err)
	}

	// FindById
	found, err := repo.FindById(ctx, id.String())
	if err != nil {
		t.Fatalf("unexpected error finding plugin: %v", err)
	}
	if found == nil || found.Name != "Test Plugin" {
		t.Errorf("expected found plugin, got %v", found)
	}

	// DeleteById
	err = repo.DeleteById(ctx, id.String())
	if err != nil {
		t.Fatalf("unexpected error deleting plugin: %v", err)
	}

	deleted, err := repo.FindById(ctx, id.String())
	if err != nil || deleted != nil {
		t.Errorf("expected nil after delete, got %v", deleted)
	}
}
