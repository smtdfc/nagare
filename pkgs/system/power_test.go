package system

import (
	"testing"
)

func TestPowerControlInstance(t *testing.T) {
	pc := NewPowerControl()
	if pc == nil {
		t.Fatalf("expected non-nil PowerControl instance")
	}

	// Test invalid action returns error
	err := pc.execute("unsupported_action_xyz")
	if err == nil {
		t.Fatalf("expected error for unsupported action, got nil")
	}
}
