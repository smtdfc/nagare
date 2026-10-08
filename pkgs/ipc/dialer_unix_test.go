//go:build !windows

package ipc

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/pkgs/paths"
)

func TestDialerUnix(t *testing.T) {
	testSocket := filepath.Join(paths.DataDir, "test_dialer.sock")
	_ = os.Remove(testSocket)
	defer os.Remove(testSocket)

	listener, err := Listen(testSocket)
	if err != nil {
		t.Fatalf("failed to create unix listener: %v", err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err == nil {
			_ = conn.Close()
		}
	}()

	ctx := context.Background()
	conn, err := Dial(ctx, testSocket)
	if err != nil {
		t.Fatalf("failed to dial unix socket: %v", err)
	}
	defer conn.Close()
}
