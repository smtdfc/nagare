package declarations

import (
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type mockExecuteBindings struct {
	mockTaskBindings
	calledName string
	success    bool
}

func (m *mockExecuteBindings) CallTool(_ *core_context.ExecuteContext, toolName string, _ string) *tool.Result {
	m.calledName = toolName
	if m.success {
		return tool.NewToolResultBuilder("call-1", toolName).Success("tool executed").Build()
	}
	return tool.NewToolResultBuilder("call-1", toolName).Failure(nil).Build()
}

func TestExecuteTool_Metadata(t *testing.T) {
	// Verify ExecuteTool metadata
	if ExecuteTool.GetName() != "execute_tool" {
		t.Errorf("expected 'execute_tool', got '%s'", ExecuteTool.GetName())
	}
	if len(ExecuteTool.GetCategories()) != 1 || ExecuteTool.GetCategories()[0] != tool.RoutingCategory {
		t.Errorf("expected ToolRoutingCategory, got %v", ExecuteTool.GetCategories())
	}
}

func TestExecuteTool_Execute(t *testing.T) {
	// Verify ExecuteTool execution success
	mb := &mockExecuteBindings{success: true}
	tl := ExecuteTool.WithBindings(mb)

	res, err := tl.Execute(nil, `{"name":"time_tool","args":{}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid output, got '%s'", res)
	}
	if mb.calledName != "time_tool" {
		t.Errorf("expected called tool 'time_tool', got '%s'", mb.calledName)
	}
}
