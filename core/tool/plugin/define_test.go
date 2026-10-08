package plugin

import (
	"testing"
)

func TestPluginTool_Metadata(t *testing.T) {
	// Verify PluginTool construction and metadata accessors
	pt := NewPluginTool(
		"custom_plugin_tool",
		`{"input":"string"}`,
		"Custom plugin tool description",
		[]string{"custom_category"},
		"plugin-123",
		nil,
	)

	if pt.GetName() != "custom_plugin_tool" {
		t.Errorf("expected name 'custom_plugin_tool', got '%s'", pt.GetName())
	}
	if pt.GetArgsSchema() != `{"input":"string"}` {
		t.Errorf("expected args schema, got '%s'", pt.GetArgsSchema())
	}
	if pt.GetDescription() != "Custom plugin tool description" {
		t.Errorf("expected description, got '%s'", pt.GetDescription())
	}
	if len(pt.GetCategories()) != 1 || pt.GetCategories()[0] != "custom_category" {
		t.Errorf("expected custom_category, got %v", pt.GetCategories())
	}
	if pt.GetBindings() != nil {
		t.Errorf("expected nil bindings")
	}

	// Verify WithBindings returns tool
	sameTool := pt.WithBindings(nil)
	if sameTool == nil {
		t.Errorf("expected non-nil tool from WithBindings")
	}
}
