package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type mockTaskBindings struct {
	createdID string
	refreshed bool
}

func (m *mockTaskBindings) RefreshTask(ctx *core_context.ExecuteContext) {
	m.refreshed = true
}

func (m *mockTaskBindings) CreateTask(ctx *core_context.ExecuteContext, sessionID string, name string, prompt string, triggerBy string, repeat bool, repeatRule string, startTime string, endTime string) (string, error) {
	m.createdID = "task-mock-123"
	return m.createdID, nil
}

func (m *mockTaskBindings) FindToolsByCategories(ctx *core_context.ExecuteContext, categories []string) ([]tool.ToolMetadata, error) {
	return nil, nil
}

func (m *mockTaskBindings) CallTool(ctx *core_context.ExecuteContext, toolName string, args string) *tool.Result {
	return nil
}

func TestCreateTaskTool_Metadata(t *testing.T) {
	// Verify CreateTaskTool metadata
	if CreateTaskTool.GetName() != "create_task_tool" {
		t.Errorf("expected 'create_task_tool', got '%s'", CreateTaskTool.GetName())
	}
	if len(CreateTaskTool.GetCategories()) != 1 || CreateTaskTool.GetCategories()[0] != tool.TaskManagementCategory {
		t.Errorf("expected TaskManagementCategory, got %v", CreateTaskTool.GetCategories())
	}
}

func TestCreateTaskTool_Execute(t *testing.T) {
	// Verify CreateTaskTool execution
	mb := &mockTaskBindings{}
	tTool := CreateTaskTool.WithBindings(mb)

	ctx := &core_context.ExecuteContext{
		Context:   context.Background(),
		SessionID: "sess-1",
	}

	res, err := tTool.Execute(ctx, `{"name":"Backup","prompt":"Do backup","repeat":false}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid task output json, got '%s'", res)
	}
	if !mb.refreshed {
		t.Errorf("expected RefreshTask to be called")
	}
}
