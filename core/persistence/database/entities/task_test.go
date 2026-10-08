package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestTaskEntity_BeforeCreate(t *testing.T) {
	// Verify BeforeCreate sets UUID when nil
	tsk := &Task{
		ID:   uuid.Nil,
		Name: "Test Task",
	}

	err := tsk.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsk.ID == uuid.Nil {
		t.Errorf("expected generated UUID, got nil")
	}

	// Verify BeforeCreate preserves existing UUID
	existingID := uuid.New()
	tsk2 := &Task{ID: existingID}
	_ = tsk2.BeforeCreate(nil)
	if tsk2.ID != existingID {
		t.Errorf("expected preserved UUID %s, got %s", existingID, tsk2.ID)
	}
}
