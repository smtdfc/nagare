//go:build !windows

package ipc

import (
	"context"
	"net"
)

// Dial connects to a Unix domain socket.
func Dial(ctx context.Context, path string) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(ctx, "unix", path)
}
