package task

import (
	"testing"
)

func TestTaskTriggerSource_ToString(t *testing.T) {
	// Verify trigger source string representation
	if Scheduled.ToString() != "scheduled" {
		t.Errorf("expected 'scheduled', got '%s'", Scheduled.ToString())
	}
	if Event.ToString() != "event" {
		t.Errorf("expected 'event', got '%s'", Event.ToString())
	}
}

func TestMapTaskTriggerSourceToTaskSource(t *testing.T) {
	// Verify mapping of string to TaskTriggerSource
	tests := []struct {
		input    string
		expected TriggerSource
	}{
		{"scheduled", Scheduled},
		{"event", Event},
		{"manual", Scheduled},
		{"", Scheduled},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			res := MapTaskTriggerSourceToTaskSource(tc.input)
			if res != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, res)
			}
		})
	}
}
