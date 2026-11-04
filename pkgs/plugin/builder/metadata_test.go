package builder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuilderMetadataLoading(t *testing.T) {
	tempDir := t.TempDir()
	metaPath := filepath.Join(tempDir, "metadata.json")

	content := `{
		"package_name": "nagare.test.plugin",
		"name": "Builder Test",
		"version": "1.0.0",
		"api_version": "1.0.0",
		"author": "tester",
		"features": ["tools"]
	}`

	if err := os.WriteFile(metaPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write mock metadata: %v", err)
	}

	// Test loading valid metadata
	meta, err := loadMetadata(metaPath)
	if err != nil {
		t.Fatalf("unexpected error loading metadata: %v", err)
	}
	if meta.PackageName != "nagare.test.plugin" {
		t.Fatalf("metadata package name mismatch: %q", meta.PackageName)
	}

	// Test validating loaded metadata
	if err := validateMetadata(meta); err != nil {
		t.Fatalf("expected validation to pass: %v", err)
	}

	// Test loading from non-existent file
	_, err = loadMetadata(filepath.Join(tempDir, "non_existent.json"))
	if err == nil {
		t.Fatalf("expected error for non-existent file, got nil")
	}

	// Test saving package metadata
	pkgDir := filepath.Join(tempDir, "pkg")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatalf("failed to create pkg dir: %v", err)
	}
	binFile := filepath.Join(pkgDir, "bin", "nagare.test.plugin")

	err = savePackageMetadata(pkgDir, binFile, *meta)
	if err != nil {
		t.Fatalf("unexpected error saving package metadata: %v", err)
	}

	savedMeta, err := loadMetadata(filepath.Join(pkgDir, "metadata.json"))
	if err != nil {
		t.Fatalf("failed to load saved package metadata: %v", err)
	}
	if len(savedMeta.Bin) == 0 {
		t.Fatalf("expected binary metadata to be populated")
	}
}
