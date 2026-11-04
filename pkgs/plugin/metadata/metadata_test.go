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

	// Test path traversal in PackageName
	invalid = valid
	invalid.PackageName = "../../evil"
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for path traversal in package_name, got nil")
	}

	invalid = valid
	invalid.PackageName = "foo/bar"
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for slash in package_name, got nil")
	}

	// Test path traversal in Bin
	invalid = valid
	invalid.Bin = PluginBinaryMetadata{
		"linux": "../../../bin/sh",
	}
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for path traversal in bin, got nil")
	}

	invalid.Bin = PluginBinaryMetadata{
		"linux": "/usr/bin/python3",
	}
	if err := invalid.Validate(); err == nil {
		t.Fatalf("expected error for absolute path in bin, got nil")
	}

	// Test valid Bin
	validWithBin := valid
	validWithBin.Bin = PluginBinaryMetadata{
		"linux":   "bin/nagare.plugin.telegram",
		"windows": "bin/nagare.plugin.telegram.exe",
	}
	if err := validWithBin.Validate(); err != nil {
		t.Fatalf("unexpected validation error with valid bin: %v", err)
	}
}
