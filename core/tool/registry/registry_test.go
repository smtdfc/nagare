package registry

import (
	"testing"
)

func TestRegistry_Init(t *testing.T) {
	// Verify the registry is populated by init()
	if Len() == 0 {
		t.Fatalf("expected non-empty registry after init")
	}

	expectedTools := []string{
		"time_tool",
		"execute_tool",
		"find_tool_by_categories",
		"create_cronjob_tool",
		"weather_tool",
		"list_process_tool",
		"create_process_tool",
		"kill_process_tool",
		"power_control_tool",
		"volume_control_tool",
		"open_browser_tool",
		"read_file_tool",
		"write_file_tool",
		"list_directory_tool",
		"delete_file_tool",
		"get_user_directories_tool",
		"search_tool",
		"fetch_content_tool",
		"get_active_display_tool",
		"screenshot_tool",
		"hash_tool",
		"archive_tool",
		"trash_tool",
		"download_file_tool",
		"http_request_tool",
		"brightness_control_tool",
	}

	for _, name := range expectedTools {
		item, exists := Get(name)
		if !exists {
			t.Errorf("expected tool '%s' to be in the registry", name)
		}
		if item.Tool == nil {
			t.Errorf("expected non-nil Tool for '%s'", name)
		}
	}
}

func TestRegisterPluginTool_ReplacesOwnTool(t *testing.T) {
	err := RegisterPluginTool("plugin_replace_tool", `{"type":"object"}`, "first", []string{"plugin_test"}, "plugin-a", true, nil)
	if err != nil {
		t.Fatalf("unexpected error registering plugin tool: %v", err)
	}
	t.Cleanup(func() { UnregisterPluginTools("plugin-a") })

	err = RegisterPluginTool("plugin_replace_tool", `{"type":"object"}`, "second", []string{"plugin_test"}, "plugin-a", true, nil)
	if err != nil {
		t.Fatalf("expected a plugin to replace its own tool, got error: %v", err)
	}

	item, exists := Get("plugin_replace_tool")
	if !exists {
		t.Fatal("expected plugin tool to be registered")
	}
	if item.Tool.GetDescription() != "second" {
		t.Errorf("expected the replacement description, got %q", item.Tool.GetDescription())
	}
}

func TestRegisterPluginTool_RejectsBuiltInName(t *testing.T) {
	err := RegisterPluginTool("time_tool", `{"type":"object"}`, "hijack", nil, "plugin-a", true, nil)
	if err == nil {
		t.Fatal("expected error when a plugin registers a built-in tool name")
	}

	item, exists := Get("time_tool")
	if !exists || item.IsPluginTool {
		t.Error("expected the built-in time_tool to stay untouched")
	}
}

func TestRegisterPluginTool_RejectsOtherPluginName(t *testing.T) {
	if err := RegisterPluginTool("plugin_owned_tool", `{"type":"object"}`, "owned", nil, "plugin-a", true, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { UnregisterPluginTools("plugin-a") })

	err := RegisterPluginTool("plugin_owned_tool", `{"type":"object"}`, "hijack", nil, "plugin-b", true, nil)
	if err == nil {
		t.Fatal("expected error when another plugin registers an existing plugin tool")
	}
}

func TestUnregisterPluginTools_OnlyRemovesOwner(t *testing.T) {
	if err := RegisterPluginTool("plugin_a_tool", `{"type":"object"}`, "a", nil, "plugin-a", true, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := RegisterPluginTool("plugin_b_tool", `{"type":"object"}`, "b", nil, "plugin-b", true, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() {
		UnregisterPluginTools("plugin-a")
		UnregisterPluginTools("plugin-b")
	})

	removed := UnregisterPluginTools("plugin-a")
	if len(removed) != 1 || removed[0] != "plugin_a_tool" {
		t.Errorf("unexpected removed tools: %v", removed)
	}

	if _, exists := Get("plugin_a_tool"); exists {
		t.Error("expected plugin-a tool to be removed")
	}
	if _, exists := Get("plugin_b_tool"); !exists {
		t.Error("expected plugin-b tool to survive")
	}
}
