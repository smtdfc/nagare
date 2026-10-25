package declarations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
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
	f := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(f, []byte("hello world"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	res, err := ReadFileTool.Execute(nil, fmt.Sprintf(`{"path":%q}`, f))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out ReadFileOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.Content != "hello world" {
		t.Errorf("expected content 'hello world', got %q", out.Content)
	}
	if out.Bytes != len("hello world") {
		t.Errorf("expected 11 bytes, got %d", out.Bytes)
	}
	if out.Path != f {
		t.Errorf("expected path %q, got %q", f, out.Path)
	}
}

func TestReadFileTool_EmptyPath(t *testing.T) {
	_, err := ReadFileTool.Execute(nil, `{"path":""}`)
	if err == nil {
		t.Fatalf("expected error for an empty path")
	}
	if !strings.Contains(err.Error(), "path cannot be empty") {
		t.Errorf("expected empty-path error, got %q", err.Error())
	}
}

func TestReadFileTool_MissingFile(t *testing.T) {
	f := filepath.Join(t.TempDir(), "missing.txt")
	_, err := ReadFileTool.Execute(nil, fmt.Sprintf(`{"path":%q}`, f))
	if err == nil {
		t.Fatalf("expected error for a missing file")
	}
	if !strings.Contains(err.Error(), "failed to stat file") {
		t.Errorf("expected stat error, got %q", err.Error())
	}
}

func TestReadFileTool_Directory(t *testing.T) {
	_, err := ReadFileTool.Execute(nil, fmt.Sprintf(`{"path":%q}`, t.TempDir()))
	if err == nil {
		t.Fatalf("expected error when reading a directory")
	}
	if !strings.Contains(err.Error(), "path is a directory") {
		t.Errorf("expected directory error, got %q", err.Error())
	}
}

func TestReadFileTool_OversizedFile(t *testing.T) {
	f := filepath.Join(t.TempDir(), "big.txt")
	content := bytes.Repeat([]byte("a"), system.MaxFileToolBytes+1)
	if err := os.WriteFile(f, content, 0644); err != nil {
		t.Fatalf("failed to create oversized file: %v", err)
	}

	_, err := ReadFileTool.Execute(nil, fmt.Sprintf(`{"path":%q}`, f))
	if err == nil {
		t.Fatalf("expected error for an oversized file")
	}
	if !strings.Contains(err.Error(), "exceeds") {
		t.Errorf("expected size-limit error, got %q", err.Error())
	}
}
