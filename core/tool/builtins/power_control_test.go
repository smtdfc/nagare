package declarations

import (
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestPowerControlTool_Metadata(t *testing.T) {
	// Verify PowerControlTool metadata
	if PowerControlTool.GetName() != "power_control_tool" {
		t.Errorf("expected 'power_control_tool', got '%s'", PowerControlTool.GetName())
	}
	if len(PowerControlTool.GetCategories()) != 1 || PowerControlTool.GetCategories()[0] != tool.PowerManagementCategory {
		t.Errorf("expected PowerManagementCategory, got %v", PowerControlTool.GetCategories())
	}
}

func TestPowerControlTool_InvalidAction(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, PowerControlTool, approvals)

	// Verify unknown action does not panic and does not ask for approval
	res, err := PowerControlTool.Execute(ctx, `{"action":"invalid_test_action"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected output, got '%s'", res)
	}
	if len(approvals.calls) != 0 {
		t.Errorf("expected no approval request for an unknown action, got %d", len(approvals.calls))
	}

	// Verify empty action returns validation error
	_, err = PowerControlTool.Execute(ctx, `{"action":""}`)
	if err == nil {
		t.Errorf("expected error for empty action")
	}
}

func TestPowerControlTool_DeniedByUser(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: false}
	withBindings(t, PowerControlTool, approvals)

	res, err := PowerControlTool.Execute(ctx, `{"action":"shutdown"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}
	if !strings.Contains(res, "denied") {
		t.Errorf("expected denial message, got '%s'", res)
	}
}

func TestPowerControlTool_FailsClosedWithoutApprovalChannel(t *testing.T) {
	ctx := approvalTestContext()
	withBindings(t, PowerControlTool, nil)

	res, err := PowerControlTool.Execute(ctx, `{"action":"restart"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res, "no approval channel") {
		t.Errorf("expected fail-closed message, got '%s'", res)
	}
}
