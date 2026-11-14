package declarations

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestGetOSInfoTool_Metadata(t *testing.T) {
	if GetOSInfoTool.GetName() != "get_os_info_tool" {
		t.Errorf("expected 'get_os_info_tool', got '%s'", GetOSInfoTool.GetName())
	}
	if len(GetOSInfoTool.GetCategories()) != 1 || GetOSInfoTool.GetCategories()[0] != tool.SystemCategory {
		t.Errorf("expected SystemCategory, got %v", GetOSInfoTool.GetCategories())
	}
}

func TestGetOSInfoTool_Execute(t *testing.T) {
	res, err := GetOSInfoTool.Execute(nil, `{}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out GetOSInfoOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.OS != runtime.GOOS {
		t.Errorf("expected os %q, got %q", runtime.GOOS, out.OS)
	}
	if out.Arch != runtime.GOARCH {
		t.Errorf("expected arch %q, got %q", runtime.GOARCH, out.Arch)
	}
	if out.CPUs <= 0 {
		t.Errorf("expected a positive CPU count, got %d", out.CPUs)
	}
}
