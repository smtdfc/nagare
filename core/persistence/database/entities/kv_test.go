package entities

import (
	"testing"
	"time"
)

func TestKVEntity_Construction(t *testing.T) {
	// Verify KV entity fields
	now := time.Now()
	kv := KV{
		Key:       "setting_key",
		Value:     "setting_value",
		Scope:     "user_scope",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if kv.Key != "setting_key" {
		t.Errorf("expected Key 'setting_key', got '%s'", kv.Key)
	}
	if kv.Value != "setting_value" {
		t.Errorf("expected Value 'setting_value', got '%s'", kv.Value)
	}
	if kv.Scope != "user_scope" {
		t.Errorf("expected Scope 'user_scope', got '%s'", kv.Scope)
	}
}
