package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestOpenBrowserTool_Metadata(t *testing.T) {
	// Verify OpenBrowserTool metadata
	if OpenBrowserTool.GetName() != "open_browser_tool" {
		t.Errorf("expected 'open_browser_tool', got '%s'", OpenBrowserTool.GetName())
	}
	if len(OpenBrowserTool.GetCategories()) != 1 || OpenBrowserTool.GetCategories()[0] != tool.BrowserCategory {
		t.Errorf("expected BrowserCategory, got %v", OpenBrowserTool.GetCategories())
	}
}

func TestOpenBrowserTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, OpenBrowserTool, approvals)

	// Verify empty URL list is rejected without asking for approval
	_, err := OpenBrowserTool.Execute(ctx, `{"urls":[]}`)
	if err == nil {
		t.Errorf("expected error for empty URL list")
	}
	if len(approvals.calls) != 0 {
		t.Errorf("expected no approval request before validation, got %d", len(approvals.calls))
	}

	// Verify blank URLs are rejected without asking for approval
	_, err = OpenBrowserTool.Execute(ctx, `{"urls":["  "]}`)
	if err == nil {
		t.Errorf("expected error for blank URL")
	}
	if len(approvals.calls) != 0 {
		t.Errorf("expected no approval request for invalid URL, got %d", len(approvals.calls))
	}
}

func TestOpenBrowserTool_DeniedByUser(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: false}
	withBindings(t, OpenBrowserTool, approvals)

	_, err := OpenBrowserTool.Execute(ctx, `{"urls":["https://example.com"]}`)
	if err == nil {
		t.Fatalf("expected error when the user denies the action")
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}
}

func TestOpenBrowserTool_FailsClosedWithoutApprovalChannel(t *testing.T) {
	ctx := approvalTestContext()
	withBindings(t, OpenBrowserTool, nil)

	_, err := OpenBrowserTool.Execute(ctx, `{"urls":["https://example.com"]}`)
	if err == nil {
		t.Fatalf("expected error when no approval bindings are configured")
	}
}
