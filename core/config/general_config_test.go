package config

import (
	"testing"
)

func TestGeneralConfig(t *testing.T) {
	// Verify GeneralConfig fields assignment
	cfg := GeneralConfig{
		DefaultLLMProvider: "openai",
		DefaultLLMModel:    "gpt-4o",
	}

	if cfg.DefaultLLMProvider != "openai" {
		t.Errorf("expected provider 'openai', got '%s'", cfg.DefaultLLMProvider)
	}
	if cfg.DefaultLLMModel != "gpt-4o" {
		t.Errorf("expected model 'gpt-4o', got '%s'", cfg.DefaultLLMModel)
	}

	// Verify zero-value
	var zeroCfg GeneralConfig
	if zeroCfg.DefaultLLMProvider != "" || zeroCfg.DefaultLLMModel != "" {
		t.Errorf("expected empty zero value")
	}
}
