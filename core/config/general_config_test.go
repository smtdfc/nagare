package config

import (
	"testing"
)

func TestGeneralConfig(t *testing.T) {
	// Verify GeneralConfig fields assignment
	cfg := GeneralConfig{
		DefaultModelProvider: "openai",
		DefaultModel:         "gpt-4o",
		SummaryModelProvider: "anthropic",
		SummaryModel:         "claude-3-5-sonnet",
	}

	if cfg.DefaultModelProvider != "openai" {
		t.Errorf("expected provider 'openai', got '%s'", cfg.DefaultModelProvider)
	}
	if cfg.DefaultModel != "gpt-4o" {
		t.Errorf("expected model 'gpt-4o', got '%s'", cfg.DefaultModel)
	}
	if cfg.SummaryModelProvider != "anthropic" {
		t.Errorf("expected summary provider 'anthropic', got '%s'", cfg.SummaryModelProvider)
	}
	if cfg.SummaryModel != "claude-3-5-sonnet" {
		t.Errorf("expected summary model 'claude-3-5-sonnet', got '%s'", cfg.SummaryModel)
	}

	// Verify zero-value
	var zeroCfg GeneralConfig
	if zeroCfg.DefaultModelProvider != "" || zeroCfg.DefaultModel != "" ||
		zeroCfg.SummaryModelProvider != "" || zeroCfg.SummaryModel != "" {
		t.Errorf("expected empty zero value")
	}
}
