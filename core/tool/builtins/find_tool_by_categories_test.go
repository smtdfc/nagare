package declarations

import (
	"encoding/json"
	"errors"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type mockFindBindings struct {
	mockTaskBindings
	gotCategories []string
	tools         []tool.Metadata
	err           error
}

func (m *mockFindBindings) FindToolsByCategories(_ *core_context.ExecuteContext, categories []string) ([]tool.Metadata, error) {
	m.gotCategories = categories
	return m.tools, m.err
}

func TestFindToolByCategories_Metadata(t *testing.T) {
	// Verify FindToolByCategories metadata
	if FindToolByCategories.GetName() != "find_tool_by_categories" {
		t.Errorf("expected 'find_tool_by_categories', got '%s'", FindToolByCategories.GetName())
	}
	if len(FindToolByCategories.GetCategories()) != 1 || FindToolByCategories.GetCategories()[0] != tool.RoutingCategory {
		t.Errorf("expected ToolRoutingCategory, got %v", FindToolByCategories.GetCategories())
	}
}

func TestFindToolByCategories_Execute(t *testing.T) {
	mb := &mockFindBindings{
		tools: []tool.Metadata{
			{Name: "time_tool", Description: "Get time"},
			{Name: "read_file_tool", Description: "Read a file"},
		},
	}
	tl := withBindings(t, FindToolByCategories, mb)

	res, err := tl.Execute(nil, `{"categories":["timing","filesystem"]}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mb.gotCategories) != 2 || mb.gotCategories[0] != "timing" || mb.gotCategories[1] != "filesystem" {
		t.Errorf("expected categories [timing filesystem], got %v", mb.gotCategories)
	}

	var out FindToolByCategoriesOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if len(out.Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(out.Tools))
	}
	if out.Tools[0].Name != "time_tool" {
		t.Errorf("expected first tool 'time_tool', got %q", out.Tools[0].Name)
	}
	if out.Tools[1].Description != "Read a file" {
		t.Errorf("expected second tool description 'Read a file', got %q", out.Tools[1].Description)
	}
}

func TestFindToolByCategories_PropagatesError(t *testing.T) {
	mb := &mockFindBindings{err: errors.New("registry unavailable")}
	tl := withBindings(t, FindToolByCategories, mb)

	_, err := tl.Execute(nil, `{"categories":["timing"]}`)
	if err == nil {
		t.Fatalf("expected binding error to be propagated")
	}
	if err.Error() != "registry unavailable" {
		t.Errorf("expected 'registry unavailable', got %q", err.Error())
	}
}
