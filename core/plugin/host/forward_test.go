package host

import (
	"testing"
)

func TestForward_Structs(t *testing.T) {
	// Verify host connection lookup map
	host := &PluginHost{
		connections: make(map[string]*PluginConnection),
	}

	host.mu.Lock()
	host.connections["plug-1"] = &PluginConnection{pluginID: "plug-1"}
	host.mu.Unlock()

	host.mu.RLock()
	_, exists := host.connections["plug-1"]
	host.mu.RUnlock()

	if !exists {
		t.Errorf("expected connection to exist")
	}
}
