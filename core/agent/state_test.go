package agent

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/messages"
)

func TestState_Operations(t *testing.T) {
	// Verify State lifecycle and methods
	st := NewAgentState()

	if st.GetLoopCounter() != 0 {
		t.Errorf("expected initial loop counter 0, got %d", st.GetLoopCounter())
	}
	if st.IsToolCall() {
		t.Errorf("expected no tool calls initially")
	}

	// Append and commit message
	msg := messages.NewTextMessage(messages.USER, "hello")
	st.AppendMessage(msg)
	if len(st.GetFullMessage()) != 1 {
		t.Errorf("expected 1 message in full message, got %d", len(st.GetFullMessage()))
	}

	st.CommitMessage()
	if len(st.CurrentMessage) != 1 {
		t.Errorf("expected 1 committed message, got %d", len(st.CurrentMessage))
	}
	if len(st.PendingMessage) != 0 {
		t.Errorf("expected 0 pending messages after commit, got %d", len(st.PendingMessage))
	}

	// Add tool call
	tc := tool.NewToolCall("c-1", "tool_a", "{}")
	st.AddToolCall(tc)
	if !st.IsToolCall() {
		t.Errorf("expected IsToolCall to be true")
	}
	st.ResetToolCall()
	if st.IsToolCall() {
		t.Errorf("expected IsToolCall to be false after reset")
	}

	// Loop counter
	st.IncreaseLoopCounter()
	if st.GetLoopCounter() != 1 {
		t.Errorf("expected loop counter 1, got %d", st.GetLoopCounter())
	}

	// Full reset
	st.Reset()
	if st.GetLoopCounter() != 0 || len(st.CurrentMessage) != 0 {
		t.Errorf("expected clean state after Reset")
	}
}
