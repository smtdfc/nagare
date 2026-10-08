package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestMessageEntity_BeforeCreate(t *testing.T) {
	// Verify BeforeCreate sets UUID when nil
	msg := &Message{
		ID:          uuid.Nil,
		MessageKind: "text",
		Content:     "{}",
	}

	err := msg.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.ID == uuid.Nil {
		t.Errorf("expected generated UUID, got nil")
	}

	// Verify BeforeCreate preserves existing UUID
	existingID := uuid.New()
	msg2 := &Message{ID: existingID}
	_ = msg2.BeforeCreate(nil)
	if msg2.ID != existingID {
		t.Errorf("expected preserved UUID %s, got %s", existingID, msg2.ID)
	}
}
