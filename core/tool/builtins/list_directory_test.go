package declarations

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("a"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	if err := os.Mkdir(filepath.Join(tmpDir, "sub"), 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	res, err := ListDirectoryTool.Execute(nil, fmt.Sprintf(`{"path":%q}`, tmpDir))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out ListDirectoryOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.Path != tmpDir {
		t.Errorf("expected path %q, got %q", tmpDir, out.Path)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d (%+v)", len(out.Entries), out.Entries)
	}

	entries := make(map[string]bool, len(out.Entries))
	for _, entry := range out.Entries {
		entries[entry.Name] = entry.IsDir
	}
	if isDir, ok := entries["file1.txt"]; !ok || isDir {
		t.Errorf("expected file1.txt to be a file, got entry=%v present=%v", isDir, ok)
	}
	if isDir, ok := entries["sub"]; !ok || !isDir {
		t.Errorf("expected sub to be a directory, got entry=%v present=%v", isDir, ok)
	}
}

func TestListDirectoryTool_EmptyPathUsesWorkingDirectory(t *testing.T) {
	res, err := ListDirectoryTool.Execute(nil, `{}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out ListDirectoryOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.Path != "" {
		t.Errorf("expected echoed path to be empty, got %q", out.Path)
	}
}

func TestListDirectoryTool_MissingPath(t *testing.T) {
	_, err := ListDirectoryTool.Execute(nil, `{"path":"/definitely/not/a/real/nagare/path"}`)
	if err == nil {
		t.Fatalf("expected error for a missing directory")
	}
	if !strings.Contains(err.Error(), "failed to list directory") {
		t.Errorf("expected list failure, got %q", err.Error())
	}
}

func TestListDirectoryTool_FilePath(t *testing.T) {
	f := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(f, []byte("a"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	_, err := ListDirectoryTool.Execute(nil, fmt.Sprintf(`{"path":%q}`, f))
	if err == nil {
		t.Fatalf("expected error when listing a file path")
	}
}
