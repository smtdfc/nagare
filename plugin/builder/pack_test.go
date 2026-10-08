package builder

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestPackAndSignPlugin(t *testing.T) {
	tempDir := t.TempDir()

	// Test prepareDirectories
	pkgDir, binFile, sigFile, err := prepareDirectories(tempDir, "nagare.plugin.sample")
	if err != nil {
		t.Fatalf("unexpected error preparing directories: %v", err)
	}
	if pkgDir == "" || binFile == "" || sigFile == "" {
		t.Fatalf("expected non-empty paths from prepareDirectories")
	}

	// Create a dummy binary file
	dummyContent := []byte("binary content for testing signing and packing")
	if err := os.WriteFile(binFile, dummyContent, 0755); err != nil {
		t.Fatalf("failed to create dummy binary file: %v", err)
	}

	// Test signBinary
	if err := signBinary(binFile, sigFile); err != nil {
		t.Fatalf("unexpected error signing binary: %v", err)
	}

	sigBytes, err := os.ReadFile(sigFile)
	if err != nil || len(sigBytes) == 0 {
		t.Fatalf("failed to read signature file or signature is empty: %v", err)
	}

	// Test packPlugin
	zipPath := filepath.Join(tempDir, "output.nagare_plugin")
	if err := packPlugin(pkgDir, zipPath); err != nil {
		t.Fatalf("unexpected error packing plugin: %v", err)
	}

	// Verify zip contents
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("failed to open generated zip file: %v", err)
	}
	defer zipReader.Close()

	if len(zipReader.File) == 0 {
		t.Fatalf("expected packed zip to contain files")
	}
}
