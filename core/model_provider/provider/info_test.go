package provider

import (
	"testing"

	"github.com/google/uuid"
)

func TestModelProviderInfo_Construction(t *testing.T) {
	// Verify ModelProviderInfo struct fields
	id := uuid.New()
	info := ModelProviderInfo{
		ID:         id,
		Name:       "OpenAI",
		Compatible: OpenAICompatible,
	}

	if info.ID != id {
		t.Errorf("expected ID %s, got %s", id, info.ID)
	}
	if info.Name != "OpenAI" {
		t.Errorf("expected name 'OpenAI', got '%s'", info.Name)
	}
	if info.Compatible != OpenAICompatible {
		t.Errorf("expected OpenAICompatible, got %s", info.Compatible)
	}
}
