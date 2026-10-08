package registry

import (
	"testing"
)

func TestRegistry_Init(t *testing.T) {
	// Verify Registry is populated by init()
	if len(Registry) == 0 {
		t.Fatalf("expected non-empty Registry after init")
	}

	expectedTools := []string{
		"time_tool",
		"execute_tool",
		"find_tool_by_categories",
		"create_task_tool",
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
	}

	for _, name := range expectedTools {
		item, exists := Registry[name]
		if !exists {
			t.Errorf("expected tool '%s' to be in Registry", name)
		}
		if item.Tool == nil {
			t.Errorf("expected non-nil Tool for '%s'", name)
		}
	}
}
