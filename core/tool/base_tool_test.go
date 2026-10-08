package tool

import (
	"testing"

	"github.com/smtdfc/nagare/core/context"
)

type sampleArgs struct {
	Query string `json:"query"`
}

type sampleOutput struct {
	Response string `json:"response"`
}

func TestDefineTool_Execution(t *testing.T) {
	// Verify DefineTool and BaseTool methods
	cb := func(ctx *context.ExecuteContext, args *sampleArgs, bindings Bindings) (*sampleOutput, error) {
		return &sampleOutput{Response: "echo: " + args.Query}, nil
	}

	tl := DefineTool("echo_tool", "echo description", cb, []string{ProcessManagementCategory})

	if tl.GetName() != "echo_tool" {
		t.Errorf("expected name 'echo_tool', got '%s'", tl.GetName())
	}
	if tl.GetDescription() != "echo description" {
		t.Errorf("expected description 'echo description', got '%s'", tl.GetDescription())
	}
	if len(tl.GetCategories()) != 1 || tl.GetCategories()[0] != ProcessManagementCategory {
		t.Errorf("expected category %s, got %v", ProcessManagementCategory, tl.GetCategories())
	}
	if tl.GetArgsSchema() == "" || tl.GetArgsSchema() == "{}" {
		t.Errorf("expected non-empty json schema, got '%s'", tl.GetArgsSchema())
	}

	// Test Execute with valid json
	res, err := tl.Execute(nil, `{"query":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid output json, got '%s'", res)
	}

	// Test Execute with invalid json
	_, err = tl.Execute(nil, `invalid-json`)
	if err == nil {
		t.Errorf("expected error for invalid json")
	}

	// Test WithBindings
	mockB := &mockBindings{}
	tl = tl.WithBindings(mockB)
	if tl.GetBindings() != mockB {
		t.Errorf("expected bindings attached")
	}
}
