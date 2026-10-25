package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestCreateProcessTool_Metadata(t *testing.T) {
	// Verify CreateProcessTool metadata
	if CreateProcessTool.GetName() != "create_process_tool" {
		t.Errorf("expected 'create_process_tool', got '%s'", CreateProcessTool.GetName())
	}
	if len(CreateProcessTool.GetCategories()) != 1 || CreateProcessTool.GetCategories()[0] != tool.ProcessManagementCategory {
		t.Errorf("expected ProcessManagementCategory, got %v", CreateProcessTool.GetCategories())
	}
}

func TestCreateProcessTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, CreateProcessTool, approvals)

	// Verify error on empty command
	_, err := CreateProcessTool.Execute(ctx, `{"command":""}`)
	if err == nil {
		t.Errorf("expected error for empty command")
	}
	if len(approvals.calls) != 0 {
		t.Errorf("expected no approval request before validation, got %d", len(approvals.calls))
	}
}

func TestCreateProcessTool_DeniedByUser(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: false}
	withBindings(t, CreateProcessTool, approvals)

	_, err := CreateProcessTool.Execute(ctx, `{"command":"echo","args":["hello"]}`)
	if err == nil {
		t.Fatalf("expected error when the user denies the action")
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}
}

func TestCreateProcessTool_FailsClosedWithoutApprovalChannel(t *testing.T) {
	ctx := approvalTestContext()
	withBindings(t, CreateProcessTool, nil)

	_, err := CreateProcessTool.Execute(ctx, `{"command":"echo"}`)
	if err == nil {
		t.Fatalf("expected error when no approval bindings are configured")
	}
}
