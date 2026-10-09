//go:build unix

package memory

import (
	"runtime"
	"strings"
	"testing"
)

// TestGetVectorModule verifies the correct library filename is returned based on OS.
func TestGetVectorModule(t *testing.T) {
	result := getVectorModule()

	switch runtime.GOOS {
	case "darwin":
		if !strings.HasSuffix(result, "libnagare_vector.dylib") {
			t.Errorf("expected dylib suffix, got %s", result)
		}
	case "linux", "freebsd", "netbsd":
		if !strings.HasSuffix(result, "libnagare_vector.so") {
			t.Errorf("expected .so suffix, got %s", result)
		}
	case "windows":
		if !strings.HasSuffix(result, "nagare_vector.dll") {
			t.Errorf("expected .dll suffix, got %s", result)
		}
	}
}

// TestVectorIndex_TypeConversion verifies VectorIndex is a uintptr alias.
func TestVectorIndex_TypeConversion(t *testing.T) {
	var idx VectorIndex = 42
	if uintptr(idx) != 42 {
		t.Errorf("expected 42, got %d", idx)
	}
}

// TestVectorMemory_Create_Success verifies Create succeeds when the native function returns a non-zero handle.
func TestVectorMemory_Create_Success(t *testing.T) {
	vm := &VectorMemory{
		vectorIndexCreate: func(dim uintptr, bits uintptr) uintptr {
			return 12345
		},
		vectorIndexLastError: func() string { return "" },
	}

	err := vm.Create(128, 8)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vm.index != VectorIndex(12345) {
		t.Errorf("expected index 12345, got %d", vm.index)
	}
}

// TestVectorMemory_Create_Failure verifies Create returns an error when the native function returns zero.
func TestVectorMemory_Create_Failure_NoErrorMessage(t *testing.T) {
	vm := &VectorMemory{
		vectorIndexCreate: func(dim uintptr, bits uintptr) uintptr {
			return 0
		},
		vectorIndexLastError: func() string { return "" },
	}

	err := vm.Create(128, 8)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to create vector index") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestVectorMemory_Create_Failure_WithErrorMessage verifies the native error message is included.
func TestVectorMemory_Create_Failure_WithErrorMessage(t *testing.T) {
	vm := &VectorMemory{
		vectorIndexCreate: func(dim uintptr, bits uintptr) uintptr {
			return 0
		},
		vectorIndexLastError: func() string { return "out of memory" },
	}

	err := vm.Create(128, 8)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "out of memory") {
		t.Errorf("expected error to contain 'out of memory', got: %v", err)
	}
}

// TestVectorMemory_Load_Success verifies Load succeeds with a non-zero handle.
func TestVectorMemory_Load_Success(t *testing.T) {
	vm := &VectorMemory{
		vectorIndexLoad: func(path string) uintptr {
			return 99999
		},
	}

	err := vm.Load("/tmp/test.index")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vm.index != VectorIndex(99999) {
		t.Errorf("expected index 99999, got %d", vm.index)
	}
}

// TestVectorMemory_Load_Failure verifies Load returns an error when the native function returns zero.
func TestVectorMemory_Load_Failure(t *testing.T) {
	vm := &VectorMemory{
		vectorIndexLoad: func(path string) uintptr {
			return 0
		},
	}

	err := vm.Load("/tmp/missing.index")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to load vector index") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestVectorMemory_Sync_Success verifies Sync succeeds when native function returns 0.
func TestVectorMemory_Sync_Success(t *testing.T) {
	vm := &VectorMemory{
		index: VectorIndex(100),
		vectorIndexSync: func(index uintptr, path string) int32 {
			return 0
		},
	}

	err := vm.Sync("/tmp/test.index")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestVectorMemory_Sync_Failure_WithErrorMessage verifies Sync includes the native error message.
func TestVectorMemory_Sync_Failure_WithErrorMessage(t *testing.T) {
	vm := &VectorMemory{
		index: VectorIndex(100),
		vectorIndexSync: func(index uintptr, path string) int32 {
			return -1
		},
		vectorIndexLastError: func() string { return "disk full" },
	}

	err := vm.Sync("/tmp/test.index")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "disk full") {
		t.Errorf("expected error to contain 'disk full', got: %v", err)
	}
}

// TestVectorMemory_Sync_Failure_NoErrorMessage verifies Sync includes the error code.
func TestVectorMemory_Sync_Failure_NoErrorMessage(t *testing.T) {
	vm := &VectorMemory{
		index: VectorIndex(100),
		vectorIndexSync: func(index uintptr, path string) int32 {
			return 42
		},
		vectorIndexLastError: func() string { return "" },
	}

	err := vm.Sync("/tmp/test.index")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "code=42") {
		t.Errorf("expected error to contain 'code=42', got: %v", err)
	}
}
