package mappers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func TestLLMProviderMapper(t *testing.T) {
	// Verify LLMProviderMapper methods
	credMapper := NewCredentialMapper()
	mapper := NewLLMProviderMapper(credMapper)

	id := uuid.New()
	credID := uuid.New()

	entity := &entities.LLMProvider{
		ID:           id,
		Name:         "OpenAI",
		Compatible:   "OpenAI",
		ApiKey:       "sk-test",
		Models:       "gpt-4o,gpt-4o-mini",
		BaseURL:      "https://api.openai.com",
		CredentialID: credID,
		Credential: &entities.Credential{
			ID:     credID,
			Name:   "openai-cred",
			ApiKey: "sk-test",
		},
	}

	domain := mapper.ToDomain(entity)
	if domain == nil {
		t.Fatalf("expected non-nil domain")
	}
	if domain.ID != id {
		t.Errorf("expected ID %s, got %s", id, domain.ID)
	}
	if len(domain.Models) != 2 {
		t.Errorf("expected 2 models, got %d", len(domain.Models))
	}
	if domain.Credential == nil || domain.Credential.Name != "openai-cred" {
		t.Errorf("expected mapped credential")
	}

	info := mapper.ToProviderInfo(entity)
	if info == nil || info.ID != id || info.Name != "OpenAI" {
		t.Errorf("expected valid provider info")
	}

	entityBack := mapper.ToEntity(domain)
	if entityBack == nil || entityBack.Models != "gpt-4o,gpt-4o-mini" {
		t.Errorf("expected joined models string, got %v", entityBack)
	}

	// Verify nil handling
	if mapper.ToDomain(nil) != nil {
		t.Errorf("expected nil domain for nil entity")
	}
}
