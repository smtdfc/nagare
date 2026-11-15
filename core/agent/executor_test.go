package agent

import (
	"context"
	"errors"
	"testing"

	context2 "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/message"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type mockCancelAdapter struct{}

func (m mockCancelAdapter) GetModels(context.Context) ([]string, error) { return nil, nil }
func (m mockCancelAdapter) Send(context.Context, string, messages.ListMessage, tool.ListTool) (message.ReadOnlyChannel, error) {
	ch := make(chan messages.Message)
	close(ch)
	return ch, nil
}
func (m mockCancelAdapter) Embed(context.Context, string, []string) ([][]float32, error) {
	return nil, nil
}

func TestExecutor_HandleError(t *testing.T) {
	// Verify HandleError mapping HTTP error codes to custom_errors
	ex := &Executor{}

	tests := []struct {
		code        string
		expectedErr *custom_errors.NagareCoreError
	}{
		{"429", custom_errors.ErrModelQuotaExceed},
		{"402", custom_errors.ErrModelPaymentRequired},
		{"500", custom_errors.ErrModelProviderAdapter},
		{"503", custom_errors.ErrModelProviderAdapter},
		{"504", custom_errors.ErrModelProviderAdapter},
	}

	for _, tc := range tests {
		t.Run(tc.code, func(t *testing.T) {
			errMsg := messages.NewResponseFailedMessage(tc.code, "error details")
			err := ex.HandleError(errMsg)
			if err.Code != tc.expectedErr.Code {
				t.Errorf("expected code %v, got %v", tc.expectedErr.Code, err.Code)
			}
			// The provider cause must survive, otherwise the UI can only show a
			// generic message.
			if err.Details != "error details" {
				t.Errorf("expected provider cause to be preserved, got %q", err.Details)
			}
		})
	}
}

func TestExecutor_HandleError_FallsBackWithoutCause(t *testing.T) {
	ex := &Executor{}

	err := ex.HandleError(messages.NewResponseFailedMessage("500", ""))
	if err.Code != custom_errors.ErrModelProviderAdapter.Code {
		t.Errorf("expected %q, got %q", custom_errors.ErrModelProviderAdapter.Code, err.Code)
	}
	if err.Details != custom_errors.ErrModelProviderAdapter.Details {
		t.Errorf("expected default details, got %q", err.Details)
	}
}

func TestExecutor_Execute_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	baseLogger, err := logger.NewBaseLogger()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	state := NewAgentState()
	userMsg := messages.NewTextMessage(messages.USER, "hello")
	state.AppendMessage(userMsg)

	ex := &Executor{
		state:  state,
		logger: baseLogger,
	}
	opts := &InvokeOption{}
	out := make(chan messages.Message)
	ectx := &context2.ExecuteContext{
		Context:        ctx,
		MessageChannel: out,
	}

	go ex.Execute(ectx, "test-model", &mockCancelAdapter{}, out, opts)

	var received []messages.Message
	for msg := range out {
		received = append(received, msg)
	}

	if len(received) == 0 {
		t.Fatalf("expected at least 1 message from output")
	}

	lastMsg := received[len(received)-1]
	completed, ok := lastMsg.(*messages.AgentCompletedMessage)
	if !ok {
		t.Fatalf("expected last message to be AgentCompletedMessage, got %T", lastMsg)
	}

	if !completed.Cancel {
		t.Errorf("expected completed.Cancel to be true")
	}

	// Verify user message is still in pending state
	if len(state.PendingMessage) != 1 {
		t.Errorf("expected pending message to still contain user message, got %d", len(state.PendingMessage))
	}
}

func newTestExecutor(t *testing.T) *Executor {
	t.Helper()

	baseLogger, err := logger.NewBaseLogger()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	return &Executor{
		toolMgr: manager.NewToolManager(baseLogger, nil, nil, nil, nil),
		logger:  baseLogger,
	}
}

func collectToolNames(tools tool.ListTool) map[string]int {
	counts := make(map[string]int, len(tools))
	for _, item := range tools {
		counts[item.GetName()]++
	}

	return counts
}

func TestExecutor_PrepareTool_DeduplicatesPromotedTools(t *testing.T) {
	ex := newTestExecutor(t)

	history := messages.ListMessage{
		messages.NewToolPromoteMessage([]string{"read_file_tool", "time_tool"}),
		messages.NewToolPromoteMessage([]string{"read_file_tool", "read_file_tool", "missing_tool"}),
	}

	counts := collectToolNames(ex.PrepareTool(history))

	for name, count := range counts {
		if count != 1 {
			t.Errorf("expected tool %q to be declared once, got %d", name, count)
		}
	}

	if counts["read_file_tool"] != 1 {
		t.Errorf("expected promoted tool 'read_file_tool' to be declared, got %d", counts["read_file_tool"])
	}

	// time_tool is already available without promotion and must not be duplicated.
	if counts["time_tool"] != 1 {
		t.Errorf("expected base tool 'time_tool' to stay declared once, got %d", counts["time_tool"])
	}

	if _, exists := counts["missing_tool"]; exists {
		t.Error("expected unknown promoted tool to be skipped")
	}

	if counts["execute_tool"] != 1 {
		t.Errorf("expected base tool 'execute_tool' to be declared, got %d", counts["execute_tool"])
	}
}

func TestExecutor_BuildToolPromoteMessage(t *testing.T) {
	ex := newTestExecutor(t)

	call := tool.NewToolCall("call-1", "find_tool_by_categories", `{"categories":["filesystem"]}`)
	result := tool.NewToolResultBuilder("call-1", "find_tool_by_categories").
		Success(`{"tools":[{"name":"read_file_tool"},{"name":"time_tool"},{"name":""}]}`).
		Build()

	promote := ex.BuildToolPromoteMessage(call, result)
	if promote == nil {
		t.Fatal("expected a promote message for a successful discovery")
	}

	expected := []string{"read_file_tool", "time_tool"}
	if len(promote.Tools) != len(expected) {
		t.Fatalf("expected tools %v, got %v", expected, promote.Tools)
	}
	for i, name := range expected {
		if promote.Tools[i] != name {
			t.Errorf("expected tool %q at %d, got %q", name, i, promote.Tools[i])
		}
	}

	failed := tool.NewToolResultBuilder("call-2", "find_tool_by_categories").
		Failure(errors.New("discovery failed")).
		Build()
	if promote := ex.BuildToolPromoteMessage(call, failed); promote != nil {
		t.Error("expected no promote message for a failed discovery")
	}

	other := tool.NewToolCall("call-3", "time_tool", `{}`)
	if promote := ex.BuildToolPromoteMessage(other, result); promote != nil {
		t.Error("expected no promote message for other tools")
	}

	empty := tool.NewToolResultBuilder("call-4", "find_tool_by_categories").
		Success(`{"tools":[]}`).
		Build()
	if promote := ex.BuildToolPromoteMessage(call, empty); promote != nil {
		t.Error("expected no promote message when discovery returns no tools")
	}
}

// captureAdapter records the tool declarations of every Send call and replays a
// canned response per call.
type captureAdapter struct {
	toolsByCall [][]string
	responses   messages.ListMessage
	calls       int
}

func (c *captureAdapter) GetModels(context.Context) ([]string, error) { return nil, nil }

func (c *captureAdapter) Send(_ context.Context, _ string, _ messages.ListMessage, tools tool.ListTool) (message.ReadOnlyChannel, error) {
	names := make([]string, 0, len(tools))
	for _, item := range tools {
		names = append(names, item.GetName())
	}
	c.toolsByCall = append(c.toolsByCall, names)

	ch := make(chan messages.Message)
	if c.calls == 0 {
		go func() {
			defer close(ch)
			for _, msg := range c.responses {
				ch <- msg
			}
		}()
	} else {
		close(ch)
	}
	c.calls++

	return ch, nil
}

func (c *captureAdapter) Embed(context.Context, string, []string) ([][]float32, error) {
	return nil, nil
}

func TestExecutor_PromotesDiscoveredToolsIntoHistory(t *testing.T) {
	ex := newTestExecutor(t)
	ex.state = NewAgentState()
	ex.state.AppendMessage(messages.NewTextMessage(messages.USER, "read a file"))

	adapter := &captureAdapter{
		responses: messages.ListMessage{
			messages.NewToolCallMessage("call-1", "find_tool_by_categories", `{"categories":["filesystem"]}`),
		},
	}

	out := make(chan messages.Message)
	go ex.Execute(&context2.ExecuteContext{Context: context.Background()}, "test-model", adapter, out, &InvokeOption{})
	for range out {
	}

	if adapter.calls < 2 {
		t.Fatalf("expected at least 2 Send calls, got %d", adapter.calls)
	}

	var promote *messages.ToolPromoteMessage
	for _, msg := range ex.state.PendingMessage {
		if candidate, ok := msg.(*messages.ToolPromoteMessage); ok {
			promote = candidate
		}
	}
	if promote == nil {
		t.Fatal("expected a ToolPromoteMessage in the pending history")
	}
	if len(promote.Tools) == 0 {
		t.Fatal("expected the promote message to carry discovered tools")
	}

	declared := adapter.toolsByCall[len(adapter.toolsByCall)-1]
	seen := make(map[string]bool, len(declared))
	for _, name := range declared {
		if seen[name] {
			t.Errorf("expected tool %q to be declared once", name)
		}
		seen[name] = true
	}

	for _, name := range promote.Tools {
		if !seen[name] {
			t.Errorf("expected promoted tool %q to be declared to the model", name)
		}
	}
}
