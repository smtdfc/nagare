package tool

import (
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/plugin"
)

func TestDynamicTool_Construction(t *testing.T) {
	// Verify DynamicTool fields
	pluginID := uuid.New()
	p := &plugin.Plugin{
		ID:   pluginID,
		Name: "test-plugin",
	}

	dt := DynamicTool{
		Name:        "dynamic_a",
		Description: "dynamic tool description",
		Args:        `{"key":"val"}`,
		PluginID:    pluginID,
		Plugin:      p,
	}

	if dt.Name != "dynamic_a" {
		t.Errorf("expected name 'dynamic_a', got '%s'", dt.Name)
	}
	if dt.PluginID != pluginID {
		t.Errorf("expected plugin ID %s, got %s", pluginID, dt.PluginID)
	}
	if dt.Plugin == nil || dt.Plugin.Name != "test-plugin" {
		t.Errorf("expected plugin attached")
	}
}

func TestToolMetadata_Construction(t *testing.T) {
	// Verify ToolMetadata struct fields
	meta := Metadata{
		Name:         "meta_tool",
		Description:  "meta desc",
		Args:         "{}",
		IsPluginTool: true,
		PluginID:     "plug-1",
	}

	if meta.Name != "meta_tool" {
		t.Errorf("expected 'meta_tool', got '%s'", meta.Name)
	}
	if !meta.IsPluginTool {
		t.Errorf("expected IsPluginTool true")
	}
	if meta.PluginID != "plug-1" {
		t.Errorf("expected 'plug-1', got '%s'", meta.PluginID)
	}
}
