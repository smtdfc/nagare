package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestInitializedPaths(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to retrieve user home directory: %v", err)
	}

	expectedDataDir := filepath.Join(home, ".nagare")
	if DataDir != expectedDataDir {
		t.Fatalf("expected DataDir %q, got %q", expectedDataDir, DataDir)
	}

	if !strings.HasPrefix(ConfigDir, DataDir) {
		t.Fatalf("expected ConfigDir to be under DataDir: %s", ConfigDir)
	}
	if !strings.HasPrefix(DatabaseDir, DataDir) {
		t.Fatalf("expected DatabaseDir to be under DataDir: %s", DatabaseDir)
	}
	if !strings.HasPrefix(PluginDir, DataDir) {
		t.Fatalf("expected PluginDir to be under DataDir: %s", PluginDir)
	}

	if runtime.GOOS == "windows" {
		expectedPipe := `\\.\pipe\nagare-plugin`
		if PluginSocketPath != expectedPipe {
			t.Fatalf("expected Windows socket %q, got %q", expectedPipe, PluginSocketPath)
		}
	} else {
		expectedSock := filepath.Join(DataDir, "nagare.sock")
		if PluginSocketPath != expectedSock {
			t.Fatalf("expected Unix socket %q, got %q", expectedSock, PluginSocketPath)
		}
	}
}
