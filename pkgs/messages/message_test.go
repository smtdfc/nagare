package messages

import (
	"testing"
)

func TestAnyMessage(t *testing.T) {
	anyMsg := AnyMessage{
		ID:       "msg-1",
		Type:     TextMessageType,
		InvokeID: "inv-any",
	}

	if anyMsg.GetMessageID() != "msg-1" {
		t.Fatalf("expected message ID 'msg-1', got %q", anyMsg.GetMessageID())
	}
	if anyMsg.GetMessageType() != TextMessageType {
		t.Fatalf("expected message type TextMessageType, got %v", anyMsg.GetMessageType())
	}
	if anyMsg.GetInvokeID() != "inv-any" {
		t.Fatalf("expected invoke ID 'inv-any', got %q", anyMsg.GetInvokeID())
	}

	anyMsg.SetInvokeID("inv-updated")
	if anyMsg.GetInvokeID() != "inv-updated" {
		t.Fatalf("expected updated invoke ID 'inv-updated', got %q", anyMsg.GetInvokeID())
	}
}
