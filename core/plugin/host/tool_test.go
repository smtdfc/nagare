package host

import (
	"testing"
)

func TestTool_ConnectionScopes(t *testing.T) {
	// Verify PluginConnection scopes for plugin_tool feature
	conn := &PluginConnection{
		pluginID:   "p-2",
		pluginName: "calendar",
		scopes:     []string{"plugin_tool"},
	}

	if len(conn.scopes) != 1 || conn.scopes[0] != "plugin_tool" {
		t.Errorf("expected scope 'plugin_tool', got %v", conn.scopes)
	}
}
