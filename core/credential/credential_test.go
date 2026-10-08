package credential

import (
	"testing"

	"github.com/google/uuid"
)

func TestCredential_Construction(t *testing.T) {
	// Verify Credential struct construction
	id := uuid.New()
	cred := Credential{
		ID:     id,
		Name:   "openai-key",
		ApiKey: "sk-test123456",
	}

	if cred.ID != id {
		t.Errorf("expected ID %s, got %s", id, cred.ID)
	}
	if cred.Name != "openai-key" {
		t.Errorf("expected name 'openai-key', got '%s'", cred.Name)
	}
	if cred.ApiKey != "sk-test123456" {
		t.Errorf("expected apiKey 'sk-test123456', got '%s'", cred.ApiKey)
	}
}
