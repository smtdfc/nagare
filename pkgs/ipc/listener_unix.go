//go:build !windows

package ipc

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// Listen listens on a Unix domain socket.
func Listen(path string) (net.Listener, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for socket: %w", err)
	}

	// Remove stale socket if it exists
	if _, err := os.Stat(path); err == nil {
		_ = os.Remove(path)
	}

	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on unix socket %s: %w", path, err)
	}

	_ = os.Chmod(path, 0700)
	return listener, nil
}
