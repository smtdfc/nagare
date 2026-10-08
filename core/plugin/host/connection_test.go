package host

import (
	"testing"
)

func TestPluginConnection_Closed(t *testing.T) {
	// Verify PluginConnection returns error when closed
	conn := &PluginConnection{
		closed: true,
	}

	err := conn.Send("test:event", nil, "req-1")
	if err == nil {
		t.Errorf("expected error when sending on closed connection")
	}
}
