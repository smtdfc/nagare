package custom_errors

import (
	"errors"
	"testing"
)

func TestNagareCoreError_Error(t *testing.T) {
	// Test Error() returns details string
	err := NewNagareCoreError("ERR_CODE", "detailed message")
	if err.Error() != "detailed message" {
		t.Fatalf("expected 'detailed message', got '%s'", err.Error())
	}
}

func TestNewNagareCoreError(t *testing.T) {
	// Test constructor fields initialization
	code := "TEST_CODE"
	details := "test error details"
	err := NewNagareCoreError(code, details)

	if err.Code != code {
		t.Errorf("expected code %s, got %s", code, err.Code)
	}
	if err.Details != details {
		t.Errorf("expected details %s, got %s", details, err.Details)
	}

	// Verify it satisfies error interface
	var errInterface error = err
	if !errors.Is(errInterface, err) {
		t.Errorf("expected error interface to wrap err")
	}
}
