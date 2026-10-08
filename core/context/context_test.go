package context

import (
	"context"
	"testing"
)

func TestExecuteContext(t *testing.T) {
	// Verify ExecuteContext initialization
	bg := context.Background()
	execCtx := ExecuteContext{
		Context:   bg,
		SessionID: "sess-123",
	}

	if execCtx.SessionID != "sess-123" {
		t.Errorf("expected session ID 'sess-123', got '%s'", execCtx.SessionID)
	}
	if execCtx.Err() != nil {
		t.Errorf("expected nil error on background context, got %v", execCtx.Err())
	}
}
