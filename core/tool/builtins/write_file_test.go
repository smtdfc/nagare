package declarations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestWriteFileTool_Metadata(t *testing.T) {
	// Verify WriteFileTool metadata
	if WriteFileTool.GetName() != "write_file_tool" {
		t.Errorf("expected 'write_file_tool', got '%s'", WriteFileTool.GetName())
	}
	if len(WriteFileTool.GetCategories()) != 1 || WriteFileTool.GetCategories()[0] != tool.FilesystemCategory {
		t.Errorf("expected FilesystemCategory, got %v", WriteFileTool.GetCategories())
	}
}

func TestWriteFileTool_Execute(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, WriteFileTool, approvals)

	// Verify validation error when confirm is false and no approval is requested
	_, err := WriteFileTool.Execute(ctx, `{"path":"/tmp/a.txt","content":"hi","confirm":false}`)
	if err == nil {
		t.Errorf("expected error when confirm is false")
	}
	if len(approvals.calls) != 0 {
		t.Errorf("expected no approval request before validation, got %d", len(approvals.calls))
	}

	// Verify successful file write after approval
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "new_file.txt")
	res, err := WriteFileTool.Execute(ctx, `{"path":"`+f+`","content":"hello world","confirm":true}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid output, got '%s'", res)
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected exactly one approval request, got %d", len(approvals.calls))
	}

	content, err := os.ReadFile(f)
	if err != nil || string(content) != "hello world" {
		t.Errorf("expected written content 'hello world', got '%s' (err: %v)", string(content), err)
	}
}

func TestWriteFileTool_DeniedByUser(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: false}
	withBindings(t, WriteFileTool, approvals)

	f := filepath.Join(t.TempDir(), "denied.txt")
	_, err := WriteFileTool.Execute(ctx, `{"path":"`+f+`","content":"nope","confirm":true}`)
	if err == nil {
		t.Fatalf("expected error when the user denies the action")
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}
	if _, statErr := os.Stat(f); !os.IsNotExist(statErr) {
		t.Errorf("expected file to not exist after denial, got err: %v", statErr)
	}
}

func TestWriteFileTool_FailsClosedWithoutApprovalChannel(t *testing.T) {
	ctx := approvalTestContext()
	withBindings(t, WriteFileTool, nil)

	f := filepath.Join(t.TempDir(), "no_bindings.txt")
	_, err := WriteFileTool.Execute(ctx, `{"path":"`+f+`","content":"nope","confirm":true}`)
	if err == nil {
		t.Fatalf("expected error when no approval bindings are configured")
	}
	if _, statErr := os.Stat(f); !os.IsNotExist(statErr) {
		t.Errorf("expected file to not exist without approval, got err: %v", statErr)
	}
}
