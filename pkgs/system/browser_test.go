package system

import (
	"runtime"
	"testing"
)

func TestOpenBrowserInvalidOS(t *testing.T) {
	// The function checks runtime.GOOS. Under current Linux/Darwin/Windows environment,
	// test URL syntax validation or invalid command gracefully.
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
		err := OpenBrowser("https://example.com")
		if err == nil {
			t.Fatalf("expected error for unsupported OS, got nil")
		}
	}
}
