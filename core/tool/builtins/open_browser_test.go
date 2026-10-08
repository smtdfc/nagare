package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestOpenBrowserTool_Metadata(t *testing.T) {
	// Verify OpenBrowserTool metadata
	if OpenBrowserTool.GetName() != "open_browser_tool" {
		t.Errorf("expected 'open_browser_tool', got '%s'", OpenBrowserTool.GetName())
	}
	if len(OpenBrowserTool.GetCategories()) != 1 || OpenBrowserTool.GetCategories()[0] != tool.BrowserCategory {
		t.Errorf("expected BrowserCategory, got %v", OpenBrowserTool.GetCategories())
	}
}
