package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestCreateProcessTool_Metadata(t *testing.T) {
	// Verify CreateProcessTool metadata
	if CreateProcessTool.GetName() != "create_process_tool" {
		t.Errorf("expected 'create_process_tool', got '%s'", CreateProcessTool.GetName())
	}
	if len(CreateProcessTool.GetCategories()) != 1 || CreateProcessTool.GetCategories()[0] != tool.ProcessManagementCategory {
		t.Errorf("expected ProcessManagementCategory, got %v", CreateProcessTool.GetCategories())
	}
}

func TestCreateProcessTool_Validation(t *testing.T) {
	// Verify error on empty command
	_, err := CreateProcessTool.Execute(nil, `{"command":""}`)
	if err == nil {
		t.Errorf("expected error for empty command")
	}
}
