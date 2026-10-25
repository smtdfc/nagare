package declarations

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type mockExecuteBindings struct {
	mockTaskBindings
	calledName string
	calledArgs string
	result     *tool.Result
}

func (m *mockExecuteBindings) CallTool(_ *core_context.ExecuteContext, toolName string, args string) *tool.Result {
	m.calledName = toolName
	m.calledArgs = args
	return m.result
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
	mb := &mockExecuteBindings{
		result: tool.NewToolResultBuilder("call-1", "time_tool").Success("tool executed").Build(),
	}
	tl := withBindings(t, ExecuteTool, mb)

	res, err := tl.Execute(nil, `{"name":"time_tool","args":{"timezone":"UTC","limit":3}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mb.calledName != "time_tool" {
		t.Errorf("expected called tool 'time_tool', got '%s'", mb.calledName)
	}

	var sent map[string]interface{}
	if err := json.Unmarshal([]byte(mb.calledArgs), &sent); err != nil {
		t.Fatalf("failed to decode forwarded args %q: %v", mb.calledArgs, err)
	}
	if sent["timezone"] != "UTC" {
		t.Errorf("expected forwarded timezone 'UTC', got %v", sent["timezone"])
	}
	if limit, ok := sent["limit"].(float64); !ok || limit != 3 {
		t.Errorf("expected forwarded limit 3, got %v", sent["limit"])
	}

	var out ExecuteToolOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.Result == nil || !out.Result.IsSuccess {
		t.Fatalf("expected successful nested result, got %+v", out.Result)
	}
	if out.Result.Result != "tool executed" {
		t.Errorf("expected nested result 'tool executed', got %q", out.Result.Result)
	}
}

func TestExecuteTool_PropagatesFailure(t *testing.T) {
	mb := &mockExecuteBindings{
		result: tool.NewToolResultBuilder("call-1", "time_tool").Failure(errors.New("boom")).Build(),
	}
	tl := withBindings(t, ExecuteTool, mb)

	_, err := tl.Execute(nil, `{"name":"time_tool","args":{}}`)
	if err == nil {
		t.Fatalf("expected failure to be propagated")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Errorf("expected error to contain 'boom', got %q", err.Error())
	}
}

func TestExecuteTool_NilResult(t *testing.T) {
	mb := &mockExecuteBindings{result: nil}
	tl := withBindings(t, ExecuteTool, mb)

	_, err := tl.Execute(nil, `{"name":"time_tool","args":{}}`)
	if err == nil {
		t.Fatalf("expected an error when the bound tool returns no result")
	}
	if !strings.Contains(err.Error(), "returned no result") {
		t.Errorf("expected nil-result error, got %q", err.Error())
	}
}
