package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestPluginEntity_BeforeCreate(t *testing.T) {
	// Verify BeforeCreate sets UUID when nil
	p := &Plugin{
		ID:          uuid.Nil,
		PackageName: "com.test",
		Name:        "Test",
	}

	err := p.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID == uuid.Nil {
		t.Errorf("expected generated UUID, got nil")
	}

	// Verify BeforeCreate preserves existing UUID
	existingID := uuid.New()
	p2 := &Plugin{ID: existingID}
	_ = p2.BeforeCreate(nil)
	if p2.ID != existingID {
		t.Errorf("expected preserved UUID %s, got %s", existingID, p2.ID)
	}
}
