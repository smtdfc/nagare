package tool

import (
	"testing"

	context2 "github.com/smtdfc/nagare/core/context"
)

type mockBindings struct{}

func (m *mockBindings) RefreshTask(*context2.ExecuteContext) {}

func (m *mockBindings) CreateTask(*context2.ExecuteContext, string, string, string, string, bool, string, string, string) (string, error) {
	return "task-1", nil
}

func (m *mockBindings) FindToolsByCategories(*context2.ExecuteContext, []string) ([]Metadata, error) {
	return []Metadata{{Name: "mock_tool"}}, nil
}

func (m *mockBindings) CallTool(_ *context2.ExecuteContext, toolName string, _ string) *Result {
	return NewToolResultBuilder("call-1", toolName).Success("ok").Build()
}

func TestBindings_Interface(t *testing.T) {
	// Verify mock struct implements Bindings interface
	var b Bindings = &mockBindings{}

	taskID, err := b.CreateTask(nil, "sess-1", "task-name", "prompt", "scheduled", false, "no_repeat", "", "")
	if err != nil || taskID != "task-1" {
		t.Errorf("expected task-1, got %s (err: %v)", taskID, err)
	}

	tools, err := b.FindToolsByCategories(nil, []string{"system"})
	if err != nil || len(tools) != 1 || tools[0].Name != "mock_tool" {
		t.Errorf("expected 1 mock tool, got %v (err: %v)", tools, err)
	}

	res := b.CallTool(nil, "mock_tool", "{}")
	if res == nil || !res.IsSuccess {
		t.Errorf("expected successful call result")
	}
}
