package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func TestSessionRepository_CRUD(t *testing.T) {
	db, l := setupTestDB(t)
	repo := &SessionRepository{db: db, logger: l}
	ctx := context.Background()

	id := uuid.New()
	s := &entities.Session{
		ID:        id,
		Title:     "Session 1",
		OwnerID:   "u-1",
		OwnerType: "user",
	}

	if err := db.Create(s).Error; err != nil {
		t.Fatalf("failed to insert session: %v", err)
	}

	// FindByID
	found, err := repo.FindByID(ctx, id.String())
	if err != nil {
		t.Fatalf("unexpected error finding session: %v", err)
	}
	if found == nil || found.Title != "Session 1" {
		t.Errorf("expected session with title 'Session 1', got %v", found)
	}

	// FindByOwnerType
	byOwner, err := repo.FindByOwnerType(ctx, "user")
	if err != nil || len(byOwner) != 1 {
		t.Errorf("expected 1 session by owner, got %d (err: %v)", len(byOwner), err)
	}
}
