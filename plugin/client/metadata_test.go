package client

import (
	"testing"
)

func TestLoadMetadata(t *testing.T) {
	p := NewPlugin()

	validRaw := `{
		"package_name": "nagare.plugin.test",
		"name": "Test Plugin",
		"version": "1.0.0",
		"api_version": "1.0.0",
		"author": "tester",
		"features": ["chat"]
	}`

	meta, err := p.LoadMetadata(validRaw)
	if err != nil {
		t.Fatalf("unexpected error loading metadata: %v", err)
	}

	if meta == nil || meta.PackageName != "nagare.plugin.test" || meta.Name != "Test Plugin" {
		t.Fatalf("metadata content mismatch: %+v", meta)
	}
	if p.Metadata == nil || p.Metadata.PackageName != "nagare.plugin.test" {
		t.Fatalf("PluginClient.Metadata was not assigned: %+v", p.Metadata)
	}

	// Test invalid JSON string
	invalidRaw := `{invalid_json`
	_, err = p.LoadMetadata(invalidRaw)
	if err == nil {
		t.Fatalf("expected error for malformed json metadata, got nil")
	}
}
