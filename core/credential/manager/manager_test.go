package manager

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/credential"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCredentialManager_CRUD(t *testing.T) {
	// Setup in-memory sqlite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	_ = db.AutoMigrate(&entities.Credential{})

	l := &logger.BaseLogger{Logger: *slog.Default()}
	repo := repositories.NewCredentialRepository(db, l)
	mapper := mappers.NewCredentialMapper()
	mgr := &CredentialManager{credentialRepo: repo, credentialMapper: mapper}
	ctx := context.Background()

	id := uuid.New()
	val := &credential.Credential{
		ID:     id,
		Name:   "openai-key",
		ApiKey: "sk-123456",
	}

	// Create
	created, err := mgr.Create(ctx, val)
	if err != nil {
		t.Fatalf("unexpected error on Create: %v", err)
	}
	if created.ID != id {
		t.Errorf("expected ID %s, got %s", id, created.ID)
	}

	// GetByID
	found, err := mgr.GetByID(ctx, id.String())
	if err != nil || found == nil {
		t.Fatalf("unexpected error on GetByID: %v", err)
	}
	if found.Name != "openai-key" {
		t.Errorf("expected name 'openai-key', got '%s'", found.Name)
	}

	// GetAll
	all, err := mgr.GetAll(ctx)
	if err != nil || len(all) != 1 {
		t.Errorf("expected 1 credential in GetAll, got %d (err: %v)", len(all), err)
	}

	// Delete
	err = mgr.Delete(ctx, id.String())
	if err != nil {
		t.Fatalf("unexpected error on Delete: %v", err)
	}

	deleted, err := mgr.GetByID(ctx, id.String())
	if err != nil || deleted != nil {
		t.Errorf("expected nil after delete, got %v", deleted)
	}
}
