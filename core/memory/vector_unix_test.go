//go:build unix

package memory

import (
	"testing"
)

// TestOpenLibrary_InvalidPath verifies openLibrary returns an error for a non-existent library.
func TestOpenLibrary_InvalidPath(t *testing.T) {
	_, err := openLibrary("/nonexistent/path/libfake.so")
	if err == nil {
		t.Fatal("expected error for invalid library path, got nil")
	}
}
