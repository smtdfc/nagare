package tool

import (
	"testing"

	context2 "github.com/smtdfc/nagare/core/context"
)

type mockBindings struct{}

// RequestToolApproval implements [Bindings].
func (m *mockBindings) RequestToolApproval(ctx *context2.ExecuteContext, title string) bool {
	panic("unimplemented")
}

func (m *mockBindings) RefreshCronJob(*context2.ExecuteContext) {}

func (m *mockBindings) CreateCronJob(*context2.ExecuteContext, string, string, string, string) (string, error) {
	return "cronjob-1", nil
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

	jobID, err := b.CreateCronJob(nil, "sess-1", "job-name", "prompt", "0 9 * * *")
	if err != nil || jobID != "cronjob-1" {
		t.Errorf("expected cronjob-1, got %s (err: %v)", jobID, err)
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
