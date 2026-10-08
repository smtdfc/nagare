package declarations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestReadFileTool_Metadata(t *testing.T) {
	// Verify ReadFileTool metadata
	if ReadFileTool.GetName() != "read_file_tool" {
		t.Errorf("expected 'read_file_tool', got '%s'", ReadFileTool.GetName())
	}
	if len(ReadFileTool.GetCategories()) != 1 || ReadFileTool.GetCategories()[0] != tool.FilesystemCategory {
		t.Errorf("expected FilesystemCategory, got %v", ReadFileTool.GetCategories())
	}
}

func TestReadFileTool_Execute(t *testing.T) {
	// Verify reading a temporary file
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(f, []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	res, err := ReadFileTool.Execute(nil, `{"path":"`+f+`"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected output, got '%s'", res)
	}
}
