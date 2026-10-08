package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func TestLLMProviderRepository_CRUD(t *testing.T) {
	db, l := setupTestDB(t)
	repo := &LLMProviderRepository{db: db, logger: l}
	ctx := context.Background()

	id := uuid.New()
	p := &entities.LLMProvider{
		ID:         id,
		Name:       "OpenAI Provider",
		Compatible: "OpenAI",
		ApiKey:     "sk-test",
		Models:     "gpt-4o",
	}

	// Create via db
	if err := db.Create(p).Error; err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// FindByID
	found, err := repo.FindByID(ctx, id.String())
	if err != nil {
		t.Fatalf("unexpected error finding provider: %v", err)
	}
	if found == nil || found.Name != "OpenAI Provider" {
		t.Errorf("expected found provider, got %v", found)
	}

	// FindAll
	all, err := repo.FindAll(ctx)
	if err != nil || len(all) != 1 {
		t.Errorf("expected 1 provider, got %d (err: %v)", len(all), err)
	}

	// DeleteByID
	err = repo.DeleteByID(ctx, id.String())
	if err != nil {
		t.Fatalf("unexpected error deleting provider: %v", err)
	}

	deleted, err := repo.FindByID(ctx, id.String())
	if err != nil || deleted != nil {
		t.Errorf("expected nil after delete, got %v", deleted)
	}
}
