package declarations

import (
	"encoding/json"
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

func decodeListProcess(t *testing.T, res string) ListProcessOutput {
	t.Helper()

	var out ListProcessOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	return out
}

func TestListProcessTool_NameQueryWithoutMatches(t *testing.T) {
	res, err := ListProcessTool.Execute(nil, `{"name_query":"nagare_no_such_process_zzz"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := decodeListProcess(t, res)
	if out.Total != 0 {
		t.Errorf("expected no matches, got %d", out.Total)
	}
	if len(out.Processes) != 0 {
		t.Errorf("expected an empty process list, got %d entries", len(out.Processes))
	}
}

func TestListProcessTool_Limit(t *testing.T) {
	res, err := ListProcessTool.Execute(nil, `{"limit":2}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := decodeListProcess(t, res)
	if out.Total > 2 {
		t.Errorf("expected at most 2 processes, got %d", out.Total)
	}
	if len(out.Processes) != out.Total {
		t.Errorf("expected Total (%d) to match len(Processes) (%d)", out.Total, len(out.Processes))
	}
}
