package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/system"
)

func TestBrightnessControlTool_Metadata(t *testing.T) {
	if BrightnessControlTool.GetName() != "brightness_control_tool" {
		t.Errorf("expected 'brightness_control_tool', got '%s'", BrightnessControlTool.GetName())
	}
	if len(BrightnessControlTool.GetCategories()) == 0 || BrightnessControlTool.GetCategories()[0] != "system" {
		t.Errorf("expected SystemCategory, got %v", BrightnessControlTool.GetCategories())
	}
}

func TestBrightnessControlTool_Validation(t *testing.T) {
	if _, err := BrightnessControlTool.Execute(nil, `{}`); err == nil {
		t.Errorf("expected error when action is missing")
	}

	res, err := BrightnessControlTool.Execute(nil, `{"action":"rotate"}`)
	if err != nil {
		t.Fatalf("unknown actions must be reported, not fatal: %v", err)
	}
	if !contains(res, "Unknown action") {
		t.Errorf("expected 'Unknown action' message in %s", res)
	}
}

func TestBrightnessControlTool_Executes(t *testing.T) {
	// The backend may be unavailable in CI containers; the tool must report
	// failure as a result instead of a hard error in that case.
	res, err := BrightnessControlTool.Execute(nil, `{"action":"get"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(res, "success") {
		t.Errorf("expected JSON output containing 'success', got %s", res)
	}

	if err := system.NewBrightnessControl().Set(200); err == nil {
		t.Errorf("expected out of range value to be rejected")
	}
}
