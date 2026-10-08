package declarations

import (
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
	// Verify execution returns valid directories
	res, err := GetUserDirectoriesTool.Execute(nil, `{}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == "" || res == "{}" {
		t.Errorf("expected valid output, got '%s'", res)
	}
}
