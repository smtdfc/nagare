package messages

import (
	"testing"
)

func TestTextMessage(t *testing.T) {
	msg := NewTextMessage(USER, "Hello world")
	if msg == nil {
		t.Fatalf("expected non-nil TextMessage")
	}

	if msg.GetMessageID() == "" {
		t.Fatalf("expected generated UUID, got empty string")
	}
	if msg.GetMessageType() != TextMessageType {
		t.Fatalf("expected TextMessageType, got %v", msg.GetMessageType())
	}
	if msg.Content != "Hello world" || msg.Role != USER {
		t.Fatalf("message fields mismatch: %+v", msg)
	}

	msg.SetInvokeID("invoke-001")
	if msg.GetInvokeID() != "invoke-001" {
		t.Fatalf("expected invoke ID 'invoke-001', got %q", msg.GetInvokeID())
	}
}
