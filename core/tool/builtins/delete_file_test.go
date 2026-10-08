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
	// Verify validation when confirm is false
	_, err := DeleteFileTool.Execute(nil, `{"path":"/tmp/fake","confirm":false}`)
	if err == nil {
		t.Errorf("expected error when confirm is false")
	}

	// Verify successful file deletion with temp file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "to_delete.txt")
	err = os.WriteFile(testFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	res, err := DeleteFileTool.Execute(nil, `{"path":"`+testFile+`","confirm":true}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid output, got '%s'", res)
	}

	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("expected file to be deleted")
	}
}
