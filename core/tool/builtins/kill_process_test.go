package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestKillProcessTool_Metadata(t *testing.T) {
	// Verify KillProcessTool metadata
	if KillProcessTool.GetName() != "kill_process_tool" {
		t.Errorf("expected 'kill_process_tool', got '%s'", KillProcessTool.GetName())
	}
	if len(KillProcessTool.GetCategories()) != 1 || KillProcessTool.GetCategories()[0] != tool.ProcessManagementCategory {
		t.Errorf("expected ProcessManagementCategory, got %v", KillProcessTool.GetCategories())
	}
}

func TestKillProcessTool_Validation(t *testing.T) {
	// Verify error when PID is invalid
	_, err := KillProcessTool.Execute(nil, `{"pid":0}`)
	if err == nil {
		t.Errorf("expected error for PID 0")
	}

	_, err = KillProcessTool.Execute(nil, `{"pid":-5}`)
	if err == nil {
		t.Errorf("expected error for negative PID")
	}
}
