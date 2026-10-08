package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestListProcessTool_Metadata(t *testing.T) {
	// Verify ListProcessTool metadata
	if ListProcessTool.GetName() != "list_process_tool" {
		t.Errorf("expected 'list_process_tool', got '%s'", ListProcessTool.GetName())
	}
	if len(ListProcessTool.GetCategories()) != 1 || ListProcessTool.GetCategories()[0] != tool.ProcessManagementCategory {
		t.Errorf("expected ProcessManagementCategory, got %v", ListProcessTool.GetCategories())
	}
}

func TestListProcessTool_Execute(t *testing.T) {
	// Verify listing processes with limit
	res, err := ListProcessTool.Execute(nil, `{"limit":5}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid output, got '%s'", res)
	}
}
