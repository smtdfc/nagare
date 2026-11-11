package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/pkgs/paths"
)

func withWebUIPIDFile(t *testing.T) string {
	t.Helper()
	original := paths.WebUIPIDFile
	t.Cleanup(func() {
		paths.WebUIPIDFile = original
	})

	file := filepath.Join(t.TempDir(), "web-ui.pid")
	paths.WebUIPIDFile = file
	return file
}

func TestTryStopWebUIWithoutDaemon(t *testing.T) {
	withWebUIPIDFile(t)

	stopped, err := TryStopWebUI()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stopped {
		t.Errorf("expected no daemon to stop")
	}
}

func TestTryStopWebUIInvalidPIDFile(t *testing.T) {
	file := withWebUIPIDFile(t)
	if err := os.WriteFile(file, []byte("not-a-pid"), 0o644); err != nil {
		t.Fatal(err)
	}

	stopped, err := TryStopWebUI()
	if err == nil {
		t.Errorf("expected an error for an invalid PID file")
	}
	if stopped {
		t.Errorf("expected no daemon to stop")
	}
	if _, statErr := os.Stat(file); !os.IsNotExist(statErr) {
		t.Errorf("expected the invalid PID file to be cleaned up, got err: %v", statErr)
	}
}

func TestTryStopWebUIStalePID(t *testing.T) {
	file := withWebUIPIDFile(t)
	if err := os.WriteFile(file, []byte("999999999"), 0o644); err != nil {
		t.Fatal(err)
	}

	stopped, err := TryStopWebUI()
	if err != nil {
		t.Fatalf("unexpected error for a stale PID: %v", err)
	}
	if stopped {
		t.Errorf("expected no daemon to stop for a stale PID")
	}
}
