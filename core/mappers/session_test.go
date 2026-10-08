package mappers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/session"
)

func TestSessionMapper(t *testing.T) {
	// Verify SessionMapper ToDomain and ToEntity
	credMapper := NewCredentialMapper()
	llmMapper := NewLLMProviderMapper(credMapper)
	mapper := NewSessionMapper(llmMapper)

	id := uuid.New()
	providerID := uuid.New()

	entity := &entities.Session{
		ID:            id,
		Title:         "Chat 1",
		ChannelID:     "ch-1",
		OwnerID:       "u-1",
		OwnerType:     "user",
		IsArchive:     false,
		CurrentModel:  "gpt-4o",
		LLMProviderID: providerID,
		LLMProvider: &entities.LLMProvider{
			ID:         providerID,
			Name:       "OpenAI",
			Compatible: "OpenAI",
		},
	}

	domain := mapper.ToDomain(entity)
	if domain == nil {
		t.Fatalf("expected non-nil domain")
	}
	if domain.ID != id {
		t.Errorf("expected ID %s, got %s", id, domain.ID)
	}
	if domain.OwnerType != session.USER {
		t.Errorf("expected owner type USER, got %s", domain.OwnerType)
	}
	if domain.LLMProvider == nil || domain.LLMProvider.Name != "OpenAI" {
		t.Errorf("expected provider OpenAI, got %v", domain.LLMProvider)
	}

	entityBack := mapper.ToEntity(domain)
	if entityBack == nil || entityBack.Title != "Chat 1" || entityBack.OwnerType != "user" {
		t.Errorf("expected mapped entity back, got %+v", entityBack)
	}

	// Verify nil handling
	if mapper.ToDomain(nil) != nil {
		t.Errorf("expected nil domain for nil entity")
	}
}
