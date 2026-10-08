package host

import (
	"testing"
)

func TestPluginHost_SocketPath(t *testing.T) {
	// Verify PluginHost socket path getter and setter
	host := &PluginHost{}

	host.SetSocketPath("/tmp/custom.sock")
	if host.GetSocketPath() != "/tmp/custom.sock" {
		t.Errorf("expected '/tmp/custom.sock', got '%s'", host.GetSocketPath())
	}
}
