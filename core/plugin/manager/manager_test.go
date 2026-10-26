package manager

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	plugin2 "github.com/smtdfc/nagare/core/mappers/plugin"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/pkgs/paths"
	"github.com/smtdfc/nagare/plugin/metadata"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createTestZip(t *testing.T, entries map[string]string) string {
	t.Helper()
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	for name, content := range entries {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("failed to create zip entry %s: %v", name, err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write content for %s: %v", name, err)
		}
	}

	if err := zw.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "plugin_test_*.zip")
	if err != nil {
		t.Fatalf("failed to create temp zip file: %v", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(buf.Bytes()); err != nil {
		t.Fatalf("failed to write temp zip content: %v", err)
	}

	return tmpFile.Name()
}

func TestUnpackPlugin_Valid(t *testing.T) {
	zipPath := createTestZip(t, map[string]string{
		"metadata.json": `{"package_name":"test.plugin"}`,
		"bin/test":      "binary content",
	})
	defer os.Remove(zipPath)

	destDir, err := os.MkdirTemp("", "plugin_dest_*")
	if err != nil {
		t.Fatalf("failed to create destDir: %v", err)
	}
	defer os.RemoveAll(destDir)

	if err := unpackPlugin(zipPath, destDir); err != nil {
		t.Fatalf("expected unpackPlugin to succeed, got %v", err)
	}

	// Verify extracted files
	metaContent, err := os.ReadFile(filepath.Join(destDir, "metadata.json"))
	if err != nil || !strings.Contains(string(metaContent), "test.plugin") {
		t.Fatalf("metadata.json not extracted properly: %v", err)
	}
}

func TestUnpackPlugin_ZipSlip_PathTraversal(t *testing.T) {
	zipPath := createTestZip(t, map[string]string{
		"../../evil.txt": "evil content",
	})
	defer os.Remove(zipPath)

	destDir, err := os.MkdirTemp("", "plugin_dest_*")
	if err != nil {
		t.Fatalf("failed to create destDir: %v", err)
	}
	defer os.RemoveAll(destDir)

	if err := unpackPlugin(zipPath, destDir); err == nil {
		t.Fatalf("expected unpackPlugin to reject path traversal, but it succeeded")
	}
}

func TestUnpackPlugin_AbsolutePaths(t *testing.T) {
	zipPath := createTestZip(t, map[string]string{
		"/etc/evil.txt": "evil content",
	})
	defer os.Remove(zipPath)

	destDir, err := os.MkdirTemp("", "plugin_dest_*")
	if err != nil {
		t.Fatalf("failed to create destDir: %v", err)
	}
	defer os.RemoveAll(destDir)

	if err := unpackPlugin(zipPath, destDir); err == nil {
		t.Fatalf("expected unpackPlugin to reject absolute paths, but it succeeded")
	}
}

func TestUnpackPlugin_RejectsSymlinks(t *testing.T) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	hdr := &zip.FileHeader{
		Name: "symlink_entry",
	}
	hdr.SetMode(os.ModeSymlink | 0777)
	w, err := zw.CreateHeader(hdr)
	if err != nil {
		t.Fatalf("failed to create zip header: %v", err)
	}
	if _, err := w.Write([]byte("/etc/passwd")); err != nil {
		t.Fatalf("failed to write symlink target: %v", err)
	}
	_ = zw.Close()

	tmpFile, err := os.CreateTemp("", "plugin_symlink_*.zip")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	_, _ = tmpFile.Write(buf.Bytes())

	destDir, err := os.MkdirTemp("", "plugin_dest_*")
	if err != nil {
		t.Fatalf("failed to create destDir: %v", err)
	}
	defer os.RemoveAll(destDir)

	if err := unpackPlugin(tmpFile.Name(), destDir); err == nil {
		t.Fatalf("expected unpackPlugin to reject symlinks, but it succeeded")
	}
}

func TestVerifyFileSignature(t *testing.T) {
	binFile, err := os.CreateTemp("", "test_bin_*")
	if err != nil {
		t.Fatalf("failed to create temp bin file: %v", err)
	}
	defer os.Remove(binFile.Name())
	defer binFile.Close()

	content := []byte("executable binary payload")
	if _, err := binFile.Write(content); err != nil {
		t.Fatalf("failed to write binary content: %v", err)
	}

	hasher := sha256.New()
	hasher.Write(content)
	validHash := hex.EncodeToString(hasher.Sum(nil))

	sigFile, err := os.CreateTemp("", "test_sig_*")
	if err != nil {
		t.Fatalf("failed to create temp sig file: %v", err)
	}
	defer os.Remove(sigFile.Name())
	defer sigFile.Close()

	if _, err := sigFile.WriteString(validHash); err != nil {
		t.Fatalf("failed to write valid sig: %v", err)
	}

	// Valid signature test
	if err := verifyFileSignature(binFile.Name(), sigFile.Name()); err != nil {
		t.Fatalf("expected signature verification to pass, got: %v", err)
	}

	// Tampered signature test
	tamperedSigFile, err := os.CreateTemp("", "test_sig_tampered_*")
	if err != nil {
		t.Fatalf("failed to create temp sig file: %v", err)
	}
	defer os.Remove(tamperedSigFile.Name())
	defer tamperedSigFile.Close()

	if _, err := tamperedSigFile.WriteString("0000000000000000000000000000000000000000000000000000000000000000"); err != nil {
		t.Fatalf("failed to write tampered sig: %v", err)
	}

	if err := verifyFileSignature(binFile.Name(), tamperedSigFile.Name()); err == nil {
		t.Fatalf("expected signature verification to fail for tampered signature, but it passed")
	}

	// Missing signature test
	if err := verifyFileSignature(binFile.Name(), "non_existent_file.sig"); err == nil {
		t.Fatalf("expected signature verification to fail for missing signature file, but it passed")
	}
}

func newTestPluginManager() *PluginManager {
	return &PluginManager{
		logger:       &logger.BaseLogger{Logger: *slog.Default()},
		connectCodes: make(map[string]string),
	}
}

func TestVerifyPluginChecksum(t *testing.T) {
	pm := newTestPluginManager()

	binPath := filepath.Join(t.TempDir(), "plugin-binary")
	if err := os.WriteFile(binPath, []byte("original binary"), 0755); err != nil {
		t.Fatalf("failed to write binary: %v", err)
	}

	checksum, err := computeFileChecksum(binPath)
	if err != nil {
		t.Fatalf("failed to compute checksum: %v", err)
	}

	cases := []struct {
		name    string
		plugin  *plugin.Plugin
		wantErr error
	}{
		{
			name:   "matching checksum",
			plugin: &plugin.Plugin{Bin: binPath, BinChecksum: checksum},
		},
		{
			name:   "matches uppercase checksum",
			plugin: &plugin.Plugin{Bin: binPath, BinChecksum: strings.ToUpper(checksum)},
		},
		{
			name:    "mismatched checksum",
			plugin:  &plugin.Plugin{Bin: binPath, BinChecksum: strings.Repeat("0", 64)},
			wantErr: custom_errors.ErrPluginChecksumMismatch,
		},
		{
			name:   "legacy plugin without checksum",
			plugin: &plugin.Plugin{Bin: binPath},
		},
		{
			name:    "missing binary",
			plugin:  &plugin.Plugin{Bin: filepath.Join(t.TempDir(), "missing"), BinChecksum: checksum},
			wantErr: custom_errors.ErrPluginChecksumMismatch,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := pm.verifyPluginChecksum(tc.plugin)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected %v, got %v", tc.wantErr, err)
			}
		})
	}

	if err := os.WriteFile(binPath, []byte("modified binary"), 0755); err != nil {
		t.Fatalf("failed to modify binary: %v", err)
	}
	if err := pm.verifyPluginChecksum(&plugin.Plugin{Bin: binPath, BinChecksum: checksum}); !errors.Is(err, custom_errors.ErrPluginChecksumMismatch) {
		t.Fatalf("expected checksum mismatch after modifying binary, got %v", err)
	}
}

func writeTrustedPublicKey(t *testing.T, dir, name string, publicKey ed25519.PublicKey) {
	t.Helper()

	der, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
	if err := os.WriteFile(filepath.Join(dir, name), pemBytes, 0644); err != nil {
		t.Fatalf("failed to write trusted public key: %v", err)
	}
}

func TestVerifyArchiveSignature(t *testing.T) {
	pm := newTestPluginManager()

	originalTrustedDir := paths.TrustedDir
	trustedDir := t.TempDir()
	paths.TrustedDir = trustedDir
	t.Cleanup(func() {
		paths.TrustedDir = originalTrustedDir
	})

	archivePath := filepath.Join(t.TempDir(), "plugin.nagare_plugin")
	if err := os.WriteFile(archivePath, []byte("archive content"), 0644); err != nil {
		t.Fatalf("failed to write archive: %v", err)
	}

	t.Run("unsigned archive", func(t *testing.T) {
		status, err := pm.verifyArchiveSignature(archivePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != plugin.SignatureStatusUnsigned {
			t.Fatalf("expected unsigned status, got %q", status)
		}
	})

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	signature, err := metadata.SignArchive(archivePath, privateKey)
	if err != nil {
		t.Fatalf("failed to sign archive: %v", err)
	}
	if err := metadata.WriteSignatureFile(archivePath+".sig", signature); err != nil {
		t.Fatalf("failed to write signature: %v", err)
	}

	t.Run("signed without trusted key", func(t *testing.T) {
		status, err := pm.verifyArchiveSignature(archivePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != plugin.SignatureStatusUnverified {
			t.Fatalf("expected unverified status, got %q", status)
		}
	})

	decoyPublicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate decoy key: %v", err)
	}
	writeTrustedPublicKey(t, trustedDir, "aaa-decoy.pem", decoyPublicKey)
	writeTrustedPublicKey(t, trustedDir, "publisher.pem", publicKey)

	t.Run("valid signature", func(t *testing.T) {
		status, err := pm.verifyArchiveSignature(archivePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != plugin.SignatureStatusVerified {
			t.Fatalf("expected verified status, got %q", status)
		}
	})

	t.Run("explicit detached signature path", func(t *testing.T) {
		signatureBytes, err := os.ReadFile(archivePath + ".sig")
		if err != nil {
			t.Fatalf("failed to read signature: %v", err)
		}

		detachedPath := filepath.Join(t.TempDir(), "uploaded-signature")
		if err := os.WriteFile(detachedPath, signatureBytes, 0644); err != nil {
			t.Fatalf("failed to write detached signature: %v", err)
		}

		status, err := pm.verifyArchiveSignature(archivePath, detachedPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != plugin.SignatureStatusVerified {
			t.Fatalf("expected verified status, got %q", status)
		}
	})

	t.Run("skips malformed trusted key files", func(t *testing.T) {
		if err := os.WriteFile(filepath.Join(trustedDir, "README.txt"), []byte("not a key"), 0644); err != nil {
			t.Fatalf("failed to write malformed key file: %v", err)
		}

		status, err := pm.verifyArchiveSignature(archivePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != plugin.SignatureStatusVerified {
			t.Fatalf("expected verified status, got %q", status)
		}
	})

	t.Run("wrong trusted key", func(t *testing.T) {
		otherPublicKey, _, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("failed to generate second key: %v", err)
		}

		previousDir := paths.TrustedDir
		paths.TrustedDir = t.TempDir()
		t.Cleanup(func() {
			paths.TrustedDir = previousDir
		})
		writeTrustedPublicKey(t, paths.TrustedDir, "wrong.pem", otherPublicKey)

		if _, err := pm.verifyArchiveSignature(archivePath); !errors.Is(err, metadata.ErrInvalidSignature) {
			t.Fatalf("expected invalid signature error, got %v", err)
		}
	})

	t.Run("tampered archive", func(t *testing.T) {
		tamperedPath := filepath.Join(t.TempDir(), "tampered.nagare_plugin")
		if err := os.WriteFile(tamperedPath, []byte("tampered content"), 0644); err != nil {
			t.Fatalf("failed to write tampered archive: %v", err)
		}
		if err := metadata.WriteSignatureFile(tamperedPath+".sig", signature); err != nil {
			t.Fatalf("failed to write signature: %v", err)
		}

		if _, err := pm.verifyArchiveSignature(tamperedPath); !errors.Is(err, metadata.ErrInvalidSignature) {
			t.Fatalf("expected invalid signature error, got %v", err)
		}
	})

	t.Run("malformed signature file", func(t *testing.T) {
		badPath := filepath.Join(t.TempDir(), "bad.nagare_plugin")
		if err := os.WriteFile(badPath, []byte("archive content"), 0644); err != nil {
			t.Fatalf("failed to write archive: %v", err)
		}
		if err := os.WriteFile(badPath+".sig", []byte("not-a-signature"), 0644); err != nil {
			t.Fatalf("failed to write signature: %v", err)
		}

		if _, err := pm.verifyArchiveSignature(badPath); err == nil {
			t.Fatalf("expected malformed signature to fail")
		}
	})
}

func TestPluginManager_ConcurrentAccess(t *testing.T) {
	pm := &PluginManager{
		connectCodes: make(map[string]string),
	}

	var wg sync.WaitGroup
	workers := 20
	iterations := 100

	for i := 0; i < workers; i++ {
		wg.Add(1)
		pkgName := fmt.Sprintf("plugin.test.%d", i%5)
		go func(pkg string) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				pm.SetConnectCode(pkg, "code-123")
				pm.SetPluginSocketPath("/tmp/test.sock")
				_ = pm.GetPluginSocketPath()

				pm.mu.RLock()
				_ = pm.connectCodes[pkg]
				pm.mu.RUnlock()
			}
		}(pkgName)
	}

	wg.Wait()
}

func TestStartAllPlugin_ContinuesWhenPluginFails(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	if err := db.AutoMigrate(&entities.Plugin{}); err != nil {
		t.Fatalf("failed to migrate plugin entity: %v", err)
	}

	log := &logger.BaseLogger{Logger: *slog.Default()}
	repo := repositories.NewPluginRepository(db, log)
	pm := &PluginManager{
		logger:       log,
		pluginRepo:   repo,
		pluginMapper: plugin2.NewPluginMapper(),
		connectCodes: make(map[string]string),
	}

	ctx := context.Background()
	binDir := t.TempDir()
	for _, packageName := range []string{"broken.one", "broken.two"} {
		if _, err := repo.CreateOrUpdate(ctx, &entities.Plugin{
			PackageName: packageName,
			Name:        packageName,
			Version:     "1.0.0",
			Bin:         filepath.Join(binDir, packageName),
			IsActive:    true,
		}); err != nil {
			t.Fatalf("failed to seed plugin %s: %v", packageName, err)
		}
	}

	if err := pm.StartAllPlugin(ctx); err != nil {
		t.Fatalf("expected gateway startup to survive broken plugins, got %v", err)
	}

	ents, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("failed to list plugins: %v", err)
	}
	if len(ents) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(ents))
	}
	for _, ent := range ents {
		if ent.LastStatus != plugin.PluginStatusFailed {
			t.Errorf("expected failed last status for %s, got %q", ent.PackageName, ent.LastStatus)
		}
		if !ent.IsActive {
			t.Errorf("expected %s to stay active for retry, got inactive", ent.PackageName)
		}
	}
}
