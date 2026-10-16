package host

import (
	"testing"
)

func TestProvider_ConnectionScopes(t *testing.T) {
	conn := &PluginConnection{
		pluginID:   "p-1",
		pluginName: "telegram",
		scopes:     []string{"model_provider"},
	}

	if len(conn.scopes) != 1 || conn.scopes[0] != "model_provider" {
		t.Errorf("expected scope 'model_provider', got %v", conn.scopes)
	}
}
