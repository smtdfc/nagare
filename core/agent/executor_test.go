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
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}
		})
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

	out := make(chan messages.Message)
	ectx := &context2.ExecuteContext{
		Context:        ctx,
		MessageChannel: out,
	}

	go ex.Execute(ectx, "test-model", &mockCancelAdapter{}, out)

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
