package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestPowerControlTool_Metadata(t *testing.T) {
	// Verify PowerControlTool metadata
	if PowerControlTool.GetName() != "power_control_tool" {
		t.Errorf("expected 'power_control_tool', got '%s'", PowerControlTool.GetName())
	}
	if len(PowerControlTool.GetCategories()) != 1 || PowerControlTool.GetCategories()[0] != tool.PowerManagementCategory {
		t.Errorf("expected PowerManagementCategory, got %v", PowerControlTool.GetCategories())
	}
}

func TestPowerControlTool_InvalidAction(t *testing.T) {
	// Verify unknown action does not panic and returns error/failure response
	res, err := PowerControlTool.Execute(nil, `{"action":"invalid_test_action"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected output, got '%s'", res)
	}

	// Verify empty action returns validation error
	_, err = PowerControlTool.Execute(nil, `{"action":""}`)
	if err == nil {
		t.Errorf("expected error for empty action")
	}
}
