package declarations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/paths"
	"github.com/smtdfc/nagare/pkgs/system"
)

func TestScreenshotTool_Metadata(t *testing.T) {
	if ScreenshotTool.GetName() != "screenshot_tool" {
		t.Errorf("expected 'screenshot_tool', got '%s'", ScreenshotTool.GetName())
	}
	if len(ScreenshotTool.GetCategories()) != 1 || ScreenshotTool.GetCategories()[0] != tool.SystemCategory {
		t.Errorf("expected SystemCategory, got %v", ScreenshotTool.GetCategories())
	}
	if schema := ScreenshotTool.GetArgsSchema(); !strings.Contains(schema, `"index"`) {
		t.Errorf("expected schema to contain index, got %s", schema)
	}
}

func TestScreenshotTool_Execute(t *testing.T) {
	displays, err := system.GetActiveDisplay()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if displays == 0 {
		t.Skip("no active display available in this environment")
	}

	originalArtifactsDir := paths.ArtifactsDir
	paths.ArtifactsDir = t.TempDir()
	t.Cleanup(func() {
		paths.ArtifactsDir = originalArtifactsDir
	})

	res, err := ScreenshotTool.Execute(nil, `{"index":0}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out ScreenshotOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.File == "" {
		t.Fatalf("expected a screenshot file path")
	}
	if filepath.Dir(out.File) != paths.ArtifactsDir {
		t.Errorf("expected file in artifacts dir %q, got %q", paths.ArtifactsDir, out.File)
	}
	if filepath.Ext(out.File) != ".png" {
		t.Errorf("expected a .png file, got %q", out.File)
	}
	if _, err := os.Stat(out.File); err != nil {
		t.Errorf("expected screenshot file to exist: %v", err)
	}
}
