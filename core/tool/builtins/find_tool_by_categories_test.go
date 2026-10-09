package declarations

import (
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type mockFindBindings struct {
	mockTaskBindings
}

func (m *mockFindBindings) FindToolsByCategories(*core_context.ExecuteContext, []string) ([]tool.Metadata, error) {
	return []tool.Metadata{
		{Name: "time_tool", Description: "Get time"},
	}, nil
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
	// Verify FindToolByCategories execution
	mb := &mockFindBindings{}
	tl := FindToolByCategories.WithBindings(mb)

	res, err := tl.Execute(nil, `{"categories":["timing"]}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid output json, got '%s'", res)
	}
}
