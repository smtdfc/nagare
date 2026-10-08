//go:build !windows

package ipc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/pkgs/paths"
)

func TestListenerUnix(t *testing.T) {
	testSocket := filepath.Join(paths.DataDir, "test_listener.sock")
	_ = os.Remove(testSocket)
	defer os.Remove(testSocket)

	// Test creation of listener on unix socket
	listener, err := Listen(testSocket)
	if err != nil {
		t.Fatalf("failed to listen on unix socket: %v", err)
	}
	defer listener.Close()

	info, err := os.Stat(testSocket)
	if err != nil {
		t.Fatalf("socket file was not created: %v", err)
	}
	// Check socket file permissions (0700)
	if info.Mode().Perm() != 0700 {
		t.Logf("socket file permission: %v", info.Mode().Perm())
	}
}
