package chat

import (
	"testing"
)

func TestAgentInvoker_NilParams(t *testing.T) {
	// Verify error when params is nil
	invoker := &AgentInvoker{}
	_, err := invoker.Invoke(nil)
	if err == nil {
		t.Errorf("expected error when invoke params is nil")
	}
}

func TestAgentInvokeParams_Construction(t *testing.T) {
	// Verify AgentInvokeParams fields
	params := &AgentInvokeParams{
		SessionID:        "sess-123",
		SenderType:       "user",
		SenderID:         "u-1",
		SendIntoEventBus: true,
		InvokeID:         "inv-1",
	}

	if params.SessionID != "sess-123" {
		t.Errorf("expected session ID 'sess-123', got '%s'", params.SessionID)
	}
	if !params.SendIntoEventBus {
		t.Errorf("expected SendIntoEventBus true")
	}
}
