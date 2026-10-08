package provider

import (
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/credential"
)

func TestLLMProviderCompatible_ToString(t *testing.T) {
	// Verify LLMProviderCompatible ToString
	if OpenAICompatible.ToString() != "OpenAI" {
		t.Errorf("expected 'OpenAI', got '%s'", OpenAICompatible.ToString())
	}
	if UnknownCompatible.ToString() != "Unknown" {
		t.Errorf("expected 'Unknown', got '%s'", UnknownCompatible.ToString())
	}
}

func TestGetCompatibleFromString(t *testing.T) {
	// Verify compatible provider parsing
	tests := []struct {
		input    string
		expected LLMProviderCompatible
	}{
		{"OpenAI", OpenAICompatible},
		{"anthropic", UnknownCompatible},
		{"unknown", UnknownCompatible},
		{"", UnknownCompatible},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			res := GetCompatibleFromString(tc.input)
			if res != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, res)
			}
		})
	}
}

func TestLLMProviderConfig_Construction(t *testing.T) {
	// Verify LLMProviderConfig fields
	id := uuid.New()
	credID := uuid.New()
	cfg := LLMProviderConfig{
		ID:           id,
		Name:         "OpenAI Prod",
		Compatible:   OpenAICompatible,
		ApiKey:       "sk-secret",
		Models:       []string{"gpt-4o", "gpt-4o-mini"},
		BaseURL:      "https://api.openai.com/v1",
		CredentialID: credID,
		Credential: &credential.Credential{
			ID:     credID,
			Name:   "openai-cred",
			ApiKey: "sk-secret",
		},
	}

	if cfg.ID != id {
		t.Errorf("expected ID %s, got %s", id, cfg.ID)
	}
	if len(cfg.Models) != 2 {
		t.Errorf("expected 2 models, got %d", len(cfg.Models))
	}
	if cfg.Credential == nil || cfg.Credential.Name != "openai-cred" {
		t.Errorf("expected valid credential attached")
	}
}
