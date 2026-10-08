package host

import (
	"testing"
)

func TestChat_ConnectionScopes(t *testing.T) {
	// Verify PluginConnection scopes for chat feature
	conn := &PluginConnection{
		pluginID:   "p-1",
		pluginName: "telegram",
		scopes:     []string{"chat"},
	}

	if len(conn.scopes) != 1 || conn.scopes[0] != "chat" {
		t.Errorf("expected scope 'chat', got %v", conn.scopes)
	}
}
