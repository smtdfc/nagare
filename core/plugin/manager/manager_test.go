package manager

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
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
