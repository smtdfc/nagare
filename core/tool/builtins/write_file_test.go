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
	// Verify validation error when confirm is false
	_, err := WriteFileTool.Execute(nil, `{"path":"/tmp/a.txt","content":"hi","confirm":false}`)
	if err == nil {
		t.Errorf("expected error when confirm is false")
	}

	// Verify successful file write
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "new_file.txt")
	res, err := WriteFileTool.Execute(nil, `{"path":"`+f+`","content":"hello world","confirm":true}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid output, got '%s'", res)
	}

	content, err := os.ReadFile(f)
	if err != nil || string(content) != "hello world" {
		t.Errorf("expected written content 'hello world', got '%s' (err: %v)", string(content), err)
	}
}
