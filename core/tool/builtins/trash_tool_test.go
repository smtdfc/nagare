package declarations

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestTrashTool_Metadata(t *testing.T) {
	if TrashTool.GetName() != "trash_tool" {
		t.Errorf("expected 'trash_tool', got '%s'", TrashTool.GetName())
	}
	if len(TrashTool.GetCategories()) == 0 || TrashTool.GetCategories()[0] != "filesystem" {
		t.Errorf("expected FilesystemCategory, got %v", TrashTool.GetCategories())
	}
}

func TestTrashTool_MoveFile(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skipf("trash backend untested for %s in unit tests", runtime.GOOS)
	}

	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, TrashTool, approvals)

	if runtime.GOOS == "linux" {
		t.Setenv("XDG_DATA_HOME", t.TempDir())
	}

	filePath := filepath.Join(t.TempDir(), "document.txt")
	if err := os.WriteFile(filePath, []byte("trash me"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := TrashTool.Execute(ctx, `{"path":"`+filePath+`","confirm":true}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}
	if !contains(string(res), "Trash") && !contains(res, "Recycle") {
		t.Logf("trash target reported: %s", res)
	}
	if _, statErr := os.Stat(filePath); !os.IsNotExist(statErr) {
		t.Errorf("expected the original path to be gone, got err: %v", statErr)
	}
}

func TestTrashTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	withBindings(t, TrashTool, &fakeApprovalBindings{approve: true})

	if _, err := TrashTool.Execute(ctx, `{"path":"/tmp/x.txt"}`); err == nil {
		t.Errorf("expected error when confirm is false")
	}
	if _, err := TrashTool.Execute(ctx, `{"confirm":true}`); err == nil {
		t.Errorf("expected error when path is missing")
	}
}

func TestTrashTool_Denied(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skipf("trash backend untested for %s in unit tests", runtime.GOOS)
	}

	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: false}
	withBindings(t, TrashTool, approvals)

	filePath := filepath.Join(t.TempDir(), "keep.txt")
	if err := os.WriteFile(filePath, []byte("keep me"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := TrashTool.Execute(ctx, `{"path":"`+filePath+`","confirm":true}`); err == nil {
		t.Fatalf("expected error when the user denies the action")
	}
	if _, statErr := os.Stat(filePath); statErr != nil {
		t.Errorf("expected the file to remain after denial, got err: %v", statErr)
	}
}
