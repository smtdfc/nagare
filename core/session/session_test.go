package session

import (
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/llm/provider"
)

func TestOwnerType_ToString(t *testing.T) {
	// Verify OwnerType string conversions
	if USER.ToString() != "user" {
		t.Errorf("expected 'user', got '%s'", USER.ToString())
	}
	if PLUGIN.ToString() != "plugin" {
		t.Errorf("expected 'plugin', got '%s'", PLUGIN.ToString())
	}
	if SYSTEM.ToString() != "system" {
		t.Errorf("expected 'system', got '%s'", SYSTEM.ToString())
	}
	if UNKNOWN.ToString() != "unknown" {
		t.Errorf("expected 'unknown', got '%s'", UNKNOWN.ToString())
	}
}

func TestGetOwnerType(t *testing.T) {
	// Test mapping from string to OwnerType
	tests := []struct {
		input    string
		expected OwnerType
	}{
		{"user", USER},
		{"plugin", PLUGIN},
		{"system", SYSTEM},
		{"other", UNKNOWN},
		{"", UNKNOWN},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := GetOwnerType(tc.input)
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestSessionInfo_Construction(t *testing.T) {
	// Verify SessionInfo struct fields
	id := uuid.New()
	providerID := uuid.New()
	info := &Info{
		ID:              id,
		Title:           "Test Session",
		OwnerID:         "user-123",
		OwnerType:       USER,
		ChannelID:       "chan-1",
		IsArchive:       false,
		CurrentLLMModel: "gpt-4o",
		LLMProviderID:   providerID,
		LLMProvider: &provider.LLMProviderInfo{
			ID:         providerID,
			Name:       "openai",
			Compatible: provider.OpenAICompatible,
		},
	}

	if info.ID != id {
		t.Errorf("expected ID %s, got %s", id, info.ID)
	}
	if info.Title != "Test Session" {
		t.Errorf("expected title 'Test Session', got '%s'", info.Title)
	}
	if info.OwnerType != USER {
		t.Errorf("expected owner type USER, got %s", info.OwnerType)
	}
	if info.LLMProvider == nil || info.LLMProvider.Name != "openai" {
		t.Errorf("expected valid LLMProvider")
	}
}
