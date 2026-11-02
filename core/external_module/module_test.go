package external_module

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/smtdfc/nagare/pkgs/paths"
)

func TestNewExternalModule(t *testing.T) {
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	module := NewExternalModule("nagare.modules.vector")

	if module.ID != "nagare.modules.vector" {
		t.Errorf("expected ID %q, got %q", "nagare.modules.vector", module.ID)
	}
	if want := filepath.Join(paths.ModuleDir, "nagare.modules.vector"+ext); module.Bin != want {
		t.Errorf("expected Bin %q, got %q", want, module.Bin)
	}
	if want := filepath.Join(paths.ModuleDataDir, module.ID); module.DataDir != want {
		t.Errorf("expected DataDir %q, got %q", want, module.DataDir)
	}
	if want := filepath.Join(paths.ModuleConfigDir, module.ID); module.ConfigDir != want {
		t.Errorf("expected ConfigDir %q, got %q", want, module.ConfigDir)
	}
	if want := module.Bin + ".pid"; module.pidFile() != want {
		t.Errorf("expected pid file %q, got %q", want, module.pidFile())
	}
}

func TestExternalModule_StartMissingBinary(t *testing.T) {
	module := ExternalModule{
		ID:  "nagare.modules.missing",
		Bin: filepath.Join(t.TempDir(), "missing"),
	}

	if err := module.Start("1234"); err == nil {
		t.Fatal("expected an error when the module binary is missing")
	}
}

func TestExternalModule_StopWithoutPIDFile(t *testing.T) {
	module := ExternalModule{
		ID:  "nagare.modules.not-running",
		Bin: filepath.Join(t.TempDir(), "module"),
	}

	if err := module.Stop(); err != nil {
		t.Fatalf("expected no error when the module is not running, got %v", err)
	}
}
