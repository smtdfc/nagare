package declarations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestListDirectoryTool_Metadata(t *testing.T) {
	// Verify ListDirectoryTool metadata
	if ListDirectoryTool.GetName() != "list_directory_tool" {
		t.Errorf("expected 'list_directory_tool', got '%s'", ListDirectoryTool.GetName())
	}
	if len(ListDirectoryTool.GetCategories()) != 1 || ListDirectoryTool.GetCategories()[0] != tool.FilesystemCategory {
		t.Errorf("expected FilesystemCategory, got %v", ListDirectoryTool.GetCategories())
	}
}

func TestListDirectoryTool_Execute(t *testing.T) {
	// Verify listing contents of a temporary directory
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("a"), 0644)

	res, err := ListDirectoryTool.Execute(nil, `{"path":"`+tmpDir+`"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid output, got '%s'", res)
	}
}
