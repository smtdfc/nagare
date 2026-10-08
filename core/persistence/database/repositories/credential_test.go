package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func TestCredentialRepository_CRUD(t *testing.T) {
	db, l := setupTestDB(t)
	repo := NewCredentialRepository(db, l)
	ctx := context.Background()

	id := uuid.New()
	cred := &entities.Credential{
		ID:     id,
		Name:   "openai",
		ApiKey: "sk-123",
	}

	// Add
	added, err := repo.Add(ctx, cred)
	if err != nil {
		t.Fatalf("unexpected error adding credential: %v", err)
	}
	if added.ID != id {
		t.Errorf("expected ID %s, got %s", id, added.ID)
	}

	// FindByID
	found, err := repo.FindByID(ctx, id.String())
	if err != nil {
		t.Fatalf("unexpected error finding credential: %v", err)
	}
	if found == nil || found.Name != "openai" {
		t.Errorf("expected found credential with name 'openai', got %v", found)
	}

	// FindAll
	all, err := repo.FindAll(ctx)
	if err != nil || len(all) != 1 {
		t.Errorf("expected 1 credential in FindAll, got %d (err: %v)", len(all), err)
	}

	// Update
	found.Name = "openai-updated"
	err = repo.Update(ctx, found)
	if err != nil {
		t.Fatalf("unexpected error updating credential: %v", err)
	}

	// DeleteByID
	err = repo.DeleteByID(ctx, id.String())
	if err != nil {
		t.Fatalf("unexpected error deleting credential: %v", err)
	}

	deleted, err := repo.FindByID(ctx, id.String())
	if err != nil || deleted != nil {
		t.Errorf("expected nil after delete, got %v (err: %v)", deleted, err)
	}
}
