package messages

import (
	"testing"
)

func TestResponseMessages(t *testing.T) {
	// Test ResponseStartedMessage
	start := NewResponseStartedMessage()
	if start == nil || start.GetMessageType() != ResponseStartedMessageType {
		t.Fatalf("unexpected ResponseStartedMessage: %+v", start)
	}
	start.SetInvokeID("inv-resp-start")
	if start.GetInvokeID() != "inv-resp-start" {
		t.Fatalf("expected invoke ID 'inv-resp-start', got %q", start.GetInvokeID())
	}

	// Test ResponseCompletedMessage
	completed := NewResponseCompletedMessage()
	if completed == nil || completed.GetMessageType() != ResponseCompletedMessageType {
		t.Fatalf("unexpected ResponseCompletedMessage: %+v", completed)
	}
	completed.SetInvokeID("inv-resp-comp")
	if completed.GetInvokeID() != "inv-resp-comp" {
		t.Fatalf("expected invoke ID 'inv-resp-comp', got %q", completed.GetInvokeID())
	}

	// Test ResponseFailedMessage
	failed := NewResponseFailedMessage("ERR_NETWORK", "timeout connecting to LLM")
	if failed == nil || failed.GetMessageType() != ResponseFailedMessageType {
		t.Fatalf("unexpected ResponseFailedMessage: %+v", failed)
	}
	if failed.Code != "ERR_NETWORK" || failed.Cause != "timeout connecting to LLM" {
		t.Fatalf("fields mismatch in failed message: %+v", failed)
	}
	failed.SetInvokeID("inv-resp-fail")
	if failed.GetInvokeID() != "inv-resp-fail" {
		t.Fatalf("expected invoke ID 'inv-resp-fail', got %q", failed.GetInvokeID())
	}
}
