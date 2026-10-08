package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func TestMessageRepository_CRUD(t *testing.T) {
	db, l := setupTestDB(t)
	repo := &MessageRepository{db: db, logger: l}
	ctx := context.Background()

	sessID := uuid.New()
	msg := &entities.Message{
		ID:          uuid.New(),
		SessionID:   sessID,
		MessageKind: "text",
		Content:     "Hello",
	}

	// CreateBatch
	err := repo.CreateBatch(ctx, []*entities.Message{msg}, 10)
	if err != nil {
		t.Fatalf("unexpected error creating batch: %v", err)
	}

	// FindBySessionID
	found, err := repo.FindBySessionID(ctx, sessID.String())
	if err != nil {
		t.Fatalf("unexpected error finding messages: %v", err)
	}
	if len(found) != 1 {
		t.Errorf("expected 1 message, got %d", len(found))
	}

	// CreateBatch empty slice returns nil
	err = repo.CreateBatch(ctx, nil, 10)
	if err != nil {
		t.Errorf("expected nil error for empty batch")
	}
}
