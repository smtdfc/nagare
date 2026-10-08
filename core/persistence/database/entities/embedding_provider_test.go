package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestEmbeddingProviderEntity_BeforeCreate(t *testing.T) {
	// Verify BeforeCreate sets UUID when nil
	prov := &EmbeddingProvider{
		ID:   uuid.Nil,
		Name: "embed-openai",
	}

	err := prov.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prov.ID == uuid.Nil {
		t.Errorf("expected generated UUID, got nil")
	}

	// Verify BeforeCreate preserves existing UUID
	existingID := uuid.New()
	prov2 := &EmbeddingProvider{ID: existingID}
	_ = prov2.BeforeCreate(nil)
	if prov2.ID != existingID {
		t.Errorf("expected preserved UUID %s, got %s", existingID, prov2.ID)
	}
}
