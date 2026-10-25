package declarations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestDeleteFileTool_Metadata(t *testing.T) {
	// Verify DeleteFileTool metadata
	if DeleteFileTool.GetName() != "delete_file_tool" {
		t.Errorf("expected 'delete_file_tool', got '%s'", DeleteFileTool.GetName())
	}
	if len(DeleteFileTool.GetCategories()) != 1 || DeleteFileTool.GetCategories()[0] != tool.FilesystemCategory {
		t.Errorf("expected FilesystemCategory, got %v", DeleteFileTool.GetCategories())
	}
}

func TestDeleteFileTool_Execute(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, DeleteFileTool, approvals)

	// Verify validation when confirm is false and no approval is requested
	_, err := DeleteFileTool.Execute(ctx, `{"path":"/tmp/fake","confirm":false}`)
	if err == nil {
		t.Errorf("expected error when confirm is false")
	}
	if len(approvals.calls) != 0 {
		t.Errorf("expected no approval request before validation, got %d", len(approvals.calls))
	}

	// Verify successful file deletion with temp file after approval
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "to_delete.txt")
	err = os.WriteFile(testFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	res, err := DeleteFileTool.Execute(ctx, `{"path":"`+testFile+`","confirm":true}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid output, got '%s'", res)
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}

	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("expected file to be deleted")
	}
}

func TestDeleteFileTool_DeniedByUser(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: false}
	withBindings(t, DeleteFileTool, approvals)

	testFile := filepath.Join(t.TempDir(), "keep_me.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, err := DeleteFileTool.Execute(ctx, `{"path":"`+testFile+`","confirm":true}`)
	if err == nil {
		t.Fatalf("expected error when the user denies the action")
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}
	if _, statErr := os.Stat(testFile); statErr != nil {
		t.Errorf("expected file to still exist after denial, got err: %v", statErr)
	}
}

func TestDeleteFileTool_FailsClosedWithoutApprovalChannel(t *testing.T) {
	ctx := approvalTestContext()
	withBindings(t, DeleteFileTool, nil)

	testFile := filepath.Join(t.TempDir(), "keep_me.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, err := DeleteFileTool.Execute(ctx, `{"path":"`+testFile+`","confirm":true}`)
	if err == nil {
		t.Fatalf("expected error when no approval bindings are configured")
	}
	if _, statErr := os.Stat(testFile); statErr != nil {
		t.Errorf("expected file to still exist without approval, got err: %v", statErr)
	}
}
