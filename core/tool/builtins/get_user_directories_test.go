package declarations

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestGetUserDirectoriesTool_Metadata(t *testing.T) {
	// Verify GetUserDirectoriesTool metadata
	if GetUserDirectoriesTool.GetName() != "get_user_directories_tool" {
		t.Errorf("expected 'get_user_directories_tool', got '%s'", GetUserDirectoriesTool.GetName())
	}
	if len(GetUserDirectoriesTool.GetCategories()) != 1 || GetUserDirectoriesTool.GetCategories()[0] != tool.FilesystemCategory {
		t.Errorf("expected FilesystemCategory, got %v", GetUserDirectoriesTool.GetCategories())
	}
}

func TestGetUserDirectoriesTool_Execute(t *testing.T) {
	res, err := GetUserDirectoriesTool.Execute(nil, `{}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out GetUserDirectoriesOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to resolve user home directory: %v", err)
	}
	if out.Directories.Home == "" {
		t.Fatalf("expected a non-empty home directory")
	}
	if out.Directories.Home != home {
		t.Errorf("expected home %q, got %q", home, out.Directories.Home)
	}
	if out.Directories.Temp != os.TempDir() {
		t.Errorf("expected temp %q, got %q", os.TempDir(), out.Directories.Temp)
	}
}
