package tool

import (
	"errors"
	"testing"
)

func TestToolResultBuilder_Success(t *testing.T) {
	// Verify ResultBuilder on success
	builder := NewToolResultBuilder("call-1", "tool_a")
	res := builder.Success("output data").Build()

	if !res.IsSuccess {
		t.Errorf("expected success")
	}
	if res.Result != "output data" {
		t.Errorf("expected 'output data', got '%s'", res.Result)
	}

	msg := res.ToMessage()
	if msg == nil {
		t.Fatalf("expected non-nil message")
	}
	if msg.CallID != "call-1" {
		t.Errorf("expected tool call ID 'call-1', got '%s'", msg.CallID)
	}
	if msg.Name != "tool_a" {
		t.Errorf("expected tool name 'tool_a', got '%s'", msg.Name)
	}
	if msg.Result != "output data" {
		t.Errorf("expected content 'output data', got '%s'", msg.Result)
	}
}

func TestToolResultBuilder_Failure(t *testing.T) {
	// Verify ResultBuilder on failure
	builder := NewToolResultBuilder("call-2", "tool_b")
	res := builder.Failure(errors.New("failed execution")).Build()

	if res.IsSuccess {
		t.Errorf("expected failure")
	}
	if res.Result != "failed execution" {
		t.Errorf("expected 'failed execution', got '%s'", res.Result)
	}

	msg := res.ToMessage()
	if msg == nil {
		t.Fatalf("expected non-nil message")
	}
	if msg.CallID != "call-2" {
		t.Errorf("expected call ID 'call-2', got '%s'", msg.CallID)
	}
	if msg.Result != "failed execution" {
		t.Errorf("expected content 'failed execution', got '%s'", msg.Result)
	}
}

func TestToolResultBuilder_SetResult(t *testing.T) {
	// Verify SetResult method
	builder := NewToolResultBuilder("call-3", "tool_c")
	res := builder.SetResult("custom result").Build()

	if res.Result != "custom result" {
		t.Errorf("expected 'custom result', got '%s'", res.Result)
	}
}
