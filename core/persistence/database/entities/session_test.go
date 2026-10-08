package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestSessionEntity_BeforeCreate(t *testing.T) {
	// Verify BeforeCreate sets UUID when nil
	s := &Session{
		ID:    uuid.Nil,
		Title: "Test Session",
	}

	err := s.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.ID == uuid.Nil {
		t.Errorf("expected generated UUID, got nil")
	}

	// Verify BeforeCreate preserves existing UUID
	existingID := uuid.New()
	s2 := &Session{ID: existingID}
	_ = s2.BeforeCreate(nil)
	if s2.ID != existingID {
		t.Errorf("expected preserved UUID %s, got %s", existingID, s2.ID)
	}
}
