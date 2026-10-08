package tool

import (
	"testing"
)

func TestToolCall_Construction(t *testing.T) {
	// Verify NewToolCall constructor
	tc := NewToolCall("call-1", "read_file", `{"path":"/tmp/a"}`)

	if tc.CallID != "call-1" {
		t.Errorf("expected call ID 'call-1', got '%s'", tc.CallID)
	}
	if tc.Name != "read_file" {
		t.Errorf("expected name 'read_file', got '%s'", tc.Name)
	}
	if tc.Args != `{"path":"/tmp/a"}` {
		t.Errorf("expected args, got '%s'", tc.Args)
	}

	// Verify ListToolCall slice
	var list ListToolCall = []*ToolCall{tc}
	if len(list) != 1 {
		t.Errorf("expected list length 1, got %d", len(list))
	}
}
