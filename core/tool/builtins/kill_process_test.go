package declarations

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestKillProcessTool_Metadata(t *testing.T) {
	// Verify KillProcessTool metadata
	if KillProcessTool.GetName() != "kill_process_tool" {
		t.Errorf("expected 'kill_process_tool', got '%s'", KillProcessTool.GetName())
	}
	if len(KillProcessTool.GetCategories()) != 1 || KillProcessTool.GetCategories()[0] != tool.ProcessManagementCategory {
		t.Errorf("expected ProcessManagementCategory, got %v", KillProcessTool.GetCategories())
	}
}

func TestKillProcessTool_InvalidPID(t *testing.T) {
	approvals := &fakeApprovalBindings{approve: true}
	tl := withBindings(t, KillProcessTool, approvals)

	for _, args := range []string{`{"pid":0}`, `{"pid":-5}`} {
		if _, err := tl.Execute(nil, args); err == nil {
			t.Errorf("expected error for %s", args)
		}
	}

	if len(approvals.calls) != 0 {
		t.Errorf("expected no approval request for an invalid PID, got %d", len(approvals.calls))
	}
}

func TestKillProcessTool_DeniedByUser(t *testing.T) {
	approvals := &fakeApprovalBindings{approve: false}
	tl := withBindings(t, KillProcessTool, approvals)

	res, err := tl.Execute(approvalTestContext(), `{"pid":2147483647}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}

	var out KillProcessOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.Success {
		t.Errorf("expected success=false after denial")
	}
	if !strings.Contains(out.Message, "does not grant permission") {
		t.Errorf("expected denial message, got %q", out.Message)
	}
}

func TestKillProcessTool_ApprovedMissingProcess(t *testing.T) {
	approvals := &fakeApprovalBindings{approve: true}
	tl := withBindings(t, KillProcessTool, approvals)

	// 2147483647 is outside the range of valid Linux PID values, so this never
	// terminates a real process.
	res, err := tl.Execute(approvalTestContext(), `{"pid":2147483647}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}

	var out KillProcessOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.Success {
		t.Errorf("expected success=false for a missing process")
	}
	if !strings.Contains(out.Message, "Failed to terminate process PID 2147483647") {
		t.Errorf("expected failure message with PID, got %q", out.Message)
	}
}
