package declarations

import (
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

func TestVolumeControlTool_Validation(t *testing.T) {
	// Verify error on empty action
	_, err := VolumeControlTool.Execute(nil, `{"action":""}`)
	if err == nil {
		t.Errorf("expected error for empty action")
	}
}
