package messages

import (
	"testing"
)

func TestReasoningMessage(t *testing.T) {
	reasoning := NewReasoningMessage("Thought process details")
	if reasoning == nil {
		t.Fatalf("expected non-nil ReasoningMessage")
	}

	if reasoning.GetMessageType() != ReasoningMessageType {
		t.Fatalf("expected ReasoningMessageType, got %v", reasoning.GetMessageType())
	}
	if reasoning.Content != "Thought process details" {
		t.Fatalf("content mismatch: %q", reasoning.Content)
	}

	reasoning.SetInvokeID("inv-reason")
	if reasoning.GetInvokeID() != "inv-reason" {
		t.Fatalf("expected invoke ID 'inv-reason', got %q", reasoning.GetInvokeID())
	}
}
