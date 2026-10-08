package mappers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/credential"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func TestCredentialMapper_ToDomain(t *testing.T) {
	// Verify ToDomain mapping
	mapper := NewCredentialMapper()
	id := uuid.New()

	entity := &entities.Credential{
		ID:     id,
		Name:   "openai-key",
		ApiKey: "sk-123456",
	}

	domain := mapper.ToDomain(entity)
	if domain == nil {
		t.Fatalf("expected non-nil domain")
	}
	if domain.ID != id {
		t.Errorf("expected ID %s, got %s", id, domain.ID)
	}
	if domain.Name != "openai-key" {
		t.Errorf("expected name 'openai-key', got '%s'", domain.Name)
	}
	if domain.ApiKey != "sk-123456" {
		t.Errorf("expected apiKey 'sk-123456', got '%s'", domain.ApiKey)
	}

	// Verify nil handling
	if mapper.ToDomain(nil) != nil {
		t.Errorf("expected nil for nil entity")
	}
}

func TestCredentialMapper_ToEntity(t *testing.T) {
	// Verify ToEntity mapping
	mapper := NewCredentialMapper()
	id := uuid.New()

	domain := &credential.Credential{
		ID:     id,
		Name:   "anthropic-key",
		ApiKey: "sk-ant-123",
	}

	entity := mapper.ToEntity(domain)
	if entity == nil {
		t.Fatalf("expected non-nil entity")
	}
	if entity.ID != id {
		t.Errorf("expected ID %s, got %s", id, entity.ID)
	}
	if entity.Name != "anthropic-key" {
		t.Errorf("expected name 'anthropic-key', got '%s'", entity.Name)
	}

	// Verify nil handling
	if mapper.ToEntity(nil) != nil {
		t.Errorf("expected nil for nil domain")
	}
}
