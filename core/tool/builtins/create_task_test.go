package declarations

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type mockTaskBindings struct {
	createdID string
	refreshed bool
	createErr error

	createCalls   int
	gotSessionID  string
	gotName       string
	gotPrompt     string
	gotTriggerBy  string
	gotRepeat     bool
	gotRepeatRule string
	gotStartTime  string
	gotEndTime    string
}

// RequestToolApproval implements [tool.Bindings].
func (m *mockTaskBindings) RequestToolApproval(*core_context.ExecuteContext, string) bool {
	return false
}

func (m *mockTaskBindings) RefreshTask(*core_context.ExecuteContext) {
	m.refreshed = true
}

func (m *mockTaskBindings) CreateTask(_ *core_context.ExecuteContext, sessionID, name, prompt, triggerBy string, repeat bool, repeatRule, startTime, endTime string) (string, error) {
	m.createCalls++
	m.gotSessionID = sessionID
	m.gotName = name
	m.gotPrompt = prompt
	m.gotTriggerBy = triggerBy
	m.gotRepeat = repeat
	m.gotRepeatRule = repeatRule
	m.gotStartTime = startTime
	m.gotEndTime = endTime

	if m.createErr != nil {
		return "", m.createErr
	}

	return m.createdID, nil
}

func (m *mockTaskBindings) FindToolsByCategories(*core_context.ExecuteContext, []string) ([]tool.Metadata, error) {
	return nil, nil
}

func (m *mockTaskBindings) CallTool(*core_context.ExecuteContext, string, string) *tool.Result {
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
	mb := &mockTaskBindings{createdID: "task-mock-123"}
	tl := withBindings(t, CreateTaskTool, mb)

	ctx := &core_context.ExecuteContext{
		Context:   context.Background(),
		SessionID: "sess-1",
	}

	res, err := tl.Execute(ctx, `{
		"name":"Backup",
		"prompt":"Do backup",
		"repeat":true,
		"repeat_rule":"daily",
		"trigger_by":"scheduled",
		"start_time":"2026-10-25 09:00:00",
		"end_time":"2026-10-26 09:00:00"
	}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out CreateTaskOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.TaskID != "task-mock-123" {
		t.Errorf("expected task id 'task-mock-123', got %q", out.TaskID)
	}

	if mb.createCalls != 1 {
		t.Errorf("expected 1 CreateTask call, got %d", mb.createCalls)
	}
	if mb.gotSessionID != "sess-1" {
		t.Errorf("expected session id 'sess-1', got %q", mb.gotSessionID)
	}
	if mb.gotName != "Backup" {
		t.Errorf("expected name 'Backup', got %q", mb.gotName)
	}
	if mb.gotPrompt != "Do backup" {
		t.Errorf("expected prompt 'Do backup', got %q", mb.gotPrompt)
	}
	if !mb.gotRepeat {
		t.Errorf("expected repeat=true to be forwarded")
	}
	if mb.gotRepeatRule != "daily" {
		t.Errorf("expected repeat_rule 'daily', got %q", mb.gotRepeatRule)
	}
	if mb.gotTriggerBy != "scheduled" {
		t.Errorf("expected trigger_by 'scheduled', got %q", mb.gotTriggerBy)
	}
	if mb.gotStartTime != "2026-10-25 09:00:00" {
		t.Errorf("expected start_time to be forwarded, got %q", mb.gotStartTime)
	}
	if mb.gotEndTime != "2026-10-26 09:00:00" {
		t.Errorf("expected end_time to be forwarded, got %q", mb.gotEndTime)
	}
	if !mb.refreshed {
		t.Errorf("expected RefreshTask to be called after a successful create")
	}
}

func TestCreateTaskTool_ErrorDoesNotRefresh(t *testing.T) {
	mb := &mockTaskBindings{createErr: errors.New("scheduler unavailable")}
	tl := withBindings(t, CreateTaskTool, mb)

	ctx := &core_context.ExecuteContext{
		Context:   context.Background(),
		SessionID: "sess-1",
	}

	_, err := tl.Execute(ctx, `{"name":"Backup","prompt":"Do backup"}`)
	if err == nil {
		t.Fatalf("expected task creation error to be propagated")
	}
	if mb.createCalls != 1 {
		t.Errorf("expected 1 CreateTask call, got %d", mb.createCalls)
	}
	if mb.refreshed {
		t.Errorf("expected RefreshTask not to be called when creation fails")
	}
}
