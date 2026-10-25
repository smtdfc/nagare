package builder

import (
	"archive/zip"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/plugin/metadata"
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
	defer func(zipReader *zip.ReadCloser) {
		_ = zipReader.Close()
	}(zipReader)

	if len(zipReader.File) == 0 {
		t.Fatalf("expected packed zip to contain files")
	}
}

func TestSignPluginArchive(t *testing.T) {
	tempDir := t.TempDir()

	archivePath := filepath.Join(tempDir, "sample.nagare_plugin")
	if err := os.WriteFile(archivePath, []byte("archive content"), 0644); err != nil {
		t.Fatalf("failed to write archive: %v", err)
	}

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	keyPath := filepath.Join(tempDir, "plugin-signing-key")
	if err := os.WriteFile(keyPath, []byte(base64.StdEncoding.EncodeToString(privateKey.Seed())), 0600); err != nil {
		t.Fatalf("failed to write private key: %v", err)
	}

	if err := signPluginArchive(archivePath, keyPath); err != nil {
		t.Fatalf("failed to sign archive: %v", err)
	}

	signature, err := metadata.ReadSignatureFile(archivePath + ".sig")
	if err != nil {
		t.Fatalf("failed to read signature: %v", err)
	}

	if err := metadata.VerifyArchive(archivePath, signature, publicKey); err != nil {
		t.Fatalf("failed to verify archive signature: %v", err)
	}

	if err := os.WriteFile(archivePath, []byte("tampered archive"), 0644); err != nil {
		t.Fatalf("failed to tamper archive: %v", err)
	}

	if err := metadata.VerifyArchive(archivePath, signature, publicKey); err == nil {
		t.Fatalf("expected tampered archive to fail verification")
	}

	invalidKeyPath := filepath.Join(tempDir, "invalid-key")
	if err := os.WriteFile(invalidKeyPath, []byte("not-a-key"), 0600); err != nil {
		t.Fatalf("failed to write invalid key: %v", err)
	}
	if err := signPluginArchive(archivePath, invalidKeyPath); err == nil {
		t.Fatalf("expected invalid key to fail signing")
	}
}
