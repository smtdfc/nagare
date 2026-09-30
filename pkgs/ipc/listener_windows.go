//go:build windows

package ipc

import (
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

// Listen listens on a Windows named pipe.
func Listen(path string) (net.Listener, error) {
	listener, err := winio.ListenPipe(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on named pipe %s: %w", path, err)
	}
	return listener, nil
}
