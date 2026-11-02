package cronjob

import (
	"testing"
)

func TestCronJobStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected Status
	}{
		{"pending", Pending},
		{"running", Running},
		{"queued", Queued},
		{"unknown", Pending},
		{"", Pending},
	}

	for _, tt := range tests {
		got := MapStringToCronJobStatus(tt.input)
		if got != tt.expected {
			t.Errorf("MapStringToCronJobStatus(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestStatusToString(t *testing.T) {
	if Pending.ToString() != "pending" {
		t.Errorf("expected 'pending', got %q", Pending.ToString())
	}
	if Running.ToString() != "running" {
		t.Errorf("expected 'running', got %q", Running.ToString())
	}
	if Queued.ToString() != "queued" {
		t.Errorf("expected 'queued', got %q", Queued.ToString())
	}
}
