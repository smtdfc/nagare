package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveFilePath(t *testing.T) {
	// Test valid paths
	resolved, err := ResolveFilePath("test/path/file.txt")
	if err != nil {
		t.Fatalf("unexpected error resolving path: %v", err)
	}
	if resolved != filepath.Clean("test/path/file.txt") {
		t.Fatalf("expected cleaned path, got %q", resolved)
	}

	// Test empty path rejection
	_, err = ResolveFilePath("")
	if err == nil {
		t.Fatalf("expected error for empty path, got nil")
	}

	_, err = ResolveFilePath("   ")
	if err == nil {
		t.Fatalf("expected error for whitespace path, got nil")
	}
}

func TestGetUserDirectories(t *testing.T) {
	dirs, err := GetUserDirectories()
	if err != nil {
		t.Fatalf("unexpected error resolving user directories: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to retrieve user home dir: %v", err)
	}

	if dirs.Home != home {
		t.Fatalf("expected home %q, got %q", home, dirs.Home)
	}
	if !strings.HasPrefix(dirs.Downloads, home) {
		t.Fatalf("expected downloads path inside home dir: %q", dirs.Downloads)
	}
	if dirs.Temp == "" {
		t.Fatalf("expected non-empty temp directory path")
	}
}

func TestXdgDirectory(t *testing.T) {
	home := "/home/testuser"

	// Test when environment variable is unset (returns fallback)
	res := xdgDirectory("NON_EXISTENT_VAR_12345", "/default/path", home)
	if res != "/default/path" {
		t.Fatalf("expected fallback path, got %q", res)
	}

	// Test when variable equals literal "$HOME"
	t.Setenv("MOCK_XDG_DIR", "$HOME")
	res = xdgDirectory("MOCK_XDG_DIR", "/default", home)
	if res != home {
		t.Fatalf("expected home path, got %q", res)
	}

	// Test when variable has "$HOME/..." prefix
	t.Setenv("MOCK_XDG_DIR", "$HOME/CustomDocs")
	res = xdgDirectory("MOCK_XDG_DIR", "/default", home)
	expected := filepath.Join(home, "CustomDocs")
	if res != expected {
		t.Fatalf("expected %q, got %q", expected, res)
	}
}
