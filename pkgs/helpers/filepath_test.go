package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsurePathExist(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "nested", "sub", "dir")

	// Verify creation of non-existent directories
	err := EnsurePathExist(targetPath)
	if err != nil {
		t.Fatalf("unexpected error creating directory path: %v", err)
	}

	info, err := os.Stat(targetPath)
	if err != nil || !info.IsDir() {
		t.Fatalf("expected path to be a directory, err=%v", err)
	}

	// Verify idempotency when directory already exists
	err = EnsurePathExist(targetPath)
	if err != nil {
		t.Fatalf("unexpected error when path already exists: %v", err)
	}
}

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "source.txt")
	dstFile := filepath.Join(tempDir, "destination.txt")

	content := []byte("Nagare test content for file copy")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	// Test successful file copy
	if err := CopyFile(srcFile, dstFile); err != nil {
		t.Fatalf("failed to copy file: %v", err)
	}

	dstData, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	if string(dstData) != string(content) {
		t.Fatalf("content mismatch: got %q, expected %q", string(dstData), string(content))
	}

	// Test copy from non-existent source
	nonExistent := filepath.Join(tempDir, "non_existent.txt")
	if err := CopyFile(nonExistent, dstFile); err == nil {
		t.Fatalf("expected error copying non-existent file, got nil")
	}
}

func TestCopyDir(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")

	// Prepare nested directories and files in src
	nestedDir := filepath.Join(srcDir, "sub")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("failed to setup src directories: %v", err)
	}

	file1 := filepath.Join(srcDir, "file1.txt")
	file2 := filepath.Join(nestedDir, "file2.txt")
	if err := os.WriteFile(file1, []byte("file1 content"), 0644); err != nil {
		t.Fatalf("failed to write file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("file2 content"), 0644); err != nil {
		t.Fatalf("failed to write file2: %v", err)
	}

	// Test directory copy
	if err := CopyDir(srcDir, destDir); err != nil {
		t.Fatalf("failed to copy directory: %v", err)
	}

	// Verify copied structure and file contents
	destFile1 := filepath.Join(destDir, "file1.txt")
	destFile2 := filepath.Join(destDir, "sub", "file2.txt")

	data1, err := os.ReadFile(destFile1)
	if err != nil || string(data1) != "file1 content" {
		t.Fatalf("failed to verify copied file1: err=%v, content=%q", err, string(data1))
	}

	data2, err := os.ReadFile(destFile2)
	if err != nil || string(data2) != "file2 content" {
		t.Fatalf("failed to verify copied file2: err=%v, content=%q", err, string(data2))
	}

	// Test copying from non-directory path
	if err := CopyDir(file1, filepath.Join(tempDir, "invalid_dest")); err == nil {
		t.Fatalf("expected error when copying from a file instead of directory, got nil")
	}
}
