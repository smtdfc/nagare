package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestTimeTool_Metadata(t *testing.T) {
	// Verify TimeTool metadata
	if TimeTool.GetName() != "time_tool" {
		t.Errorf("expected name 'time_tool', got '%s'", TimeTool.GetName())
	}
	if len(TimeTool.GetCategories()) != 1 || TimeTool.GetCategories()[0] != tool.TimingCategory {
		t.Errorf("expected TimingCategory, got %v", TimeTool.GetCategories())
	}
	if TimeTool.GetArgsSchema() == "" {
		t.Errorf("expected non-empty schema")
	}
}

func TestTimeTool_Execute(t *testing.T) {
	// Verify TimeTool execution with empty args
	res, err := TimeTool.Execute(nil, `{}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid time output, got '%s'", res)
	}

	// Verify TimeTool execution with UTC timezone
	resUTC, err := TimeTool.Execute(nil, `{"timezone":"UTC"}`)
	if err != nil {
		t.Fatalf("unexpected error with UTC: %v", err)
	}
	if resUTC == "" || resUTC == "{}" {
		t.Errorf("expected valid UTC time output, got '%s'", resUTC)
	}
}
