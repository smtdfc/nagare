package declarations

import (
	"encoding/json"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestGetActiveDisplayTool_Metadata(t *testing.T) {
	if GetActiveDisplayTool.GetName() != "get_active_display_tool" {
		t.Errorf("expected 'get_active_display_tool', got '%s'", GetActiveDisplayTool.GetName())
	}
	if len(GetActiveDisplayTool.GetCategories()) != 1 || GetActiveDisplayTool.GetCategories()[0] != tool.SystemCategory {
		t.Errorf("expected SystemCategory, got %v", GetActiveDisplayTool.GetCategories())
	}
}

func TestGetActiveDisplayTool_Execute(t *testing.T) {
	res, err := GetActiveDisplayTool.Execute(nil, `{}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out GetActiveDisplayOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.Num < 0 {
		t.Errorf("expected a non-negative display count, got %d", out.Num)
	}
}
