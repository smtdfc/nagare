package config

import (
	"testing"
)

func TestGeneralConfig(t *testing.T) {
	// Verify GeneralConfig fields assignment
	cfg := GeneralConfig{
		DefaultModelProvider: "openai",
		DefaultModel:    "gpt-4o",
	}

	if cfg.DefaultModelProvider != "openai" {
		t.Errorf("expected provider 'openai', got '%s'", cfg.DefaultModelProvider)
	}
	if cfg.DefaultModel != "gpt-4o" {
		t.Errorf("expected model 'gpt-4o', got '%s'", cfg.DefaultModel)
	}

	// Verify zero-value
	var zeroCfg GeneralConfig
	if zeroCfg.DefaultModelProvider != "" || zeroCfg.DefaultModel != "" {
		t.Errorf("expected empty zero value")
	}
}
