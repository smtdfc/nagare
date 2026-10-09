package task

import (
	"testing"
)

func TestTaskStatus_ToString(t *testing.T) {
	// Verify TaskStatus ToString method
	if Pending.ToString() != "pending" {
		t.Errorf("expected 'pending', got '%s'", Pending.ToString())
	}
	if Running.ToString() != "running" {
		t.Errorf("expected 'running', got '%s'", Running.ToString())
	}
	if Queued.ToString() != "queued" {
		t.Errorf("expected 'queued', got '%s'", Queued.ToString())
	}
}

func TestMapStringToTaskStatus(t *testing.T) {
	// Verify string to TaskStatus mapping
	tests := []struct {
		input    string
		expected Status
	}{
		{"pending", Pending},
		{"running", Running},
		{"queued", Pending}, // default fallback in implementation
		{"unknown", Pending},
		{"", Pending},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			res := MapStringToTaskStatus(tc.input)
			if res != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, res)
			}
		})
	}
}
