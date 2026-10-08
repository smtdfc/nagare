package metadata

import (
	"testing"
)

func TestPluginMetadataValidate(t *testing.T) {
	// Test valid metadata
	valid := PluginMetadata{
		PackageName: "nagare.plugin.telegram",
		Name:        "Telegram Integration",
		Version:     "1.0.0",
		ApiVersion:  "1.0.0",
		Author:      "smtdfc",
		Features:    PluginFeatures{"chat", "notifications"},
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("unexpected validation error on valid metadata: %v", err)
	}

	if valid.GetFeatures() != "chat,notifications" {
		t.Fatalf("expected features joined as 'chat,notifications', got %q", valid.GetFeatures())
	}

	// Test missing PackageName
	invalid := valid
	invalid.PackageName = ""
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for empty package_name, got nil")
	}

	// Test missing Name
	invalid = valid
	invalid.Name = "   "
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for whitespace name, got nil")
	}

	// Test missing Version
	invalid = valid
	invalid.Version = ""
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for empty version, got nil")
	}

	// Test missing ApiVersion
	invalid = valid
	invalid.ApiVersion = ""
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for empty api_version, got nil")
	}

	// Test missing Author
	invalid = valid
	invalid.Author = ""
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for empty author, got nil")
	}
}
