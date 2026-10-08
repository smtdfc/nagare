package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestCredentialEntity_BeforeCreate(t *testing.T) {
	// Verify BeforeCreate sets UUID when nil
	cred := &Credential{
		ID:     uuid.Nil,
		Name:   "openai",
		ApiKey: "sk-test",
	}

	err := cred.BeforeCreate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.ID == uuid.Nil {
		t.Errorf("expected generated UUID, got nil")
	}

	// Verify BeforeCreate preserves existing UUID
	existingID := uuid.New()
	cred2 := &Credential{ID: existingID}
	_ = cred2.BeforeCreate(nil)
	if cred2.ID != existingID {
		t.Errorf("expected preserved UUID %s, got %s", existingID, cred2.ID)
	}
}
