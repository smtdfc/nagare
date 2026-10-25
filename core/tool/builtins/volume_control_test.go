package declarations

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestVolumeControlTool_Metadata(t *testing.T) {
	// Verify VolumeControlTool metadata
	if VolumeControlTool.GetName() != "volume_control_tool" {
		t.Errorf("expected 'volume_control_tool', got '%s'", VolumeControlTool.GetName())
	}
	if len(VolumeControlTool.GetCategories()) != 1 || VolumeControlTool.GetCategories()[0] != tool.AudioManagementCategory {
		t.Errorf("expected AudioManagementCategory, got %v", VolumeControlTool.GetCategories())
	}
}

func TestVolumeControlTool_EmptyAction(t *testing.T) {
	// Verify error on empty action without touching the system mixer
	_, err := VolumeControlTool.Execute(nil, `{"action":""}`)
	if err == nil {
		t.Errorf("expected error for empty action")
	}
}

func TestVolumeControlTool_UnknownAction(t *testing.T) {
	// Verify unknown actions are rejected before calling the system mixer
	res, err := VolumeControlTool.Execute(nil, `{"action":"bogus_action"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out VolumeControlOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.Success {
		t.Errorf("expected success=false for an unknown action")
	}
	if !strings.Contains(out.Message, "Unknown action") {
		t.Errorf("expected unknown-action message, got %q", out.Message)
	}
}
