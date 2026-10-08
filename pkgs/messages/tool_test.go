package messages

import (
	"testing"
)

func TestToolCallAndResultMessage(t *testing.T) {
	// Test ToolCallMessage
	call := NewToolCallMessage("call-1", "weather_tool", `{"lat":10}`)
	if call == nil {
		t.Fatalf("expected non-nil ToolCallMessage")
	}
	if call.GetMessageType() != ToolCallMessageType {
		t.Fatalf("expected ToolCallMessageType, got %v", call.GetMessageType())
	}
	if call.CallID != "call-1" || call.Name != "weather_tool" || call.Args != `{"lat":10}` {
		t.Fatalf("tool call fields mismatch: %+v", call)
	}

	call.SetInvokeID("inv-1")
	if call.GetInvokeID() != "inv-1" {
		t.Fatalf("expected invoke ID 'inv-1', got %q", call.GetInvokeID())
	}

	// Test ToolResultMessage
	res := NewToolResultMessage("call-1", "weather_tool", `{"temp":25}`)
	if res == nil {
		t.Fatalf("expected non-nil ToolResultMessage")
	}
	if res.GetMessageType() != ToolResultMessageType {
		t.Fatalf("expected ToolResultMessageType, got %v", res.GetMessageType())
	}
	if res.CallID != "call-1" || res.Name != "weather_tool" || res.Result != `{"temp":25}` {
		t.Fatalf("tool result fields mismatch: %+v", res)
	}

	res.SetInvokeID("inv-2")
	if res.GetInvokeID() != "inv-2" {
		t.Fatalf("expected invoke ID 'inv-2', got %q", res.GetInvokeID())
	}
}
