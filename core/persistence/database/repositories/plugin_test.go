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

func TestPluginRepository_CreateOrUpdatePersistsIntegrityFields(t *testing.T) {
	db, l := setupTestDB(t)
	repo := &PluginRepository{db: db, logger: l}
	ctx := context.Background()

	initial := &entities.Plugin{
		ID:              uuid.New(),
		PackageName:     "com.test.plugin",
		Name:            "Test Plugin",
		Version:         "1.0.0",
		BinChecksum:     "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		SignatureStatus: "unsigned",
		IsActive:        true,
	}
	if _, err := repo.CreateOrUpdate(ctx, initial); err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	updated := &entities.Plugin{
		ID:              uuid.New(),
		PackageName:     "com.test.plugin",
		Name:            "Test Plugin",
		Version:         "2.0.0",
		BinChecksum:     "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		SignatureStatus: "verified",
		IsActive:        true,
	}
	if _, err := repo.CreateOrUpdate(ctx, updated); err != nil {
		t.Fatalf("failed to update plugin: %v", err)
	}

	found, err := repo.FindByPackageName(ctx, "com.test.plugin")
	if err != nil {
		t.Fatalf("failed to find plugin: %v", err)
	}
	if found == nil {
		t.Fatalf("expected plugin to exist")
	}
	if found.BinChecksum != updated.BinChecksum {
		t.Fatalf("expected checksum %q, got %q", updated.BinChecksum, found.BinChecksum)
	}
	if found.SignatureStatus != updated.SignatureStatus {
		t.Fatalf("expected signature status %q, got %q", updated.SignatureStatus, found.SignatureStatus)
	}
	if found.Version != "2.0.0" {
		t.Fatalf("expected version 2.0.0, got %q", found.Version)
	}
}
