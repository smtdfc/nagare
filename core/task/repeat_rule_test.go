package task

import (
	"testing"
)

func TestTaskRepeatRule_ToString(t *testing.T) {
	// Verify repeat rule string conversion
	if NoRepeat.ToString() != "no_repeat" {
		t.Errorf("expected 'no_repeat', got '%s'", NoRepeat.ToString())
	}
	if Daily.ToString() != "daily" {
		t.Errorf("expected 'daily', got '%s'", Daily.ToString())
	}
	if Hourly.ToString() != "hourly" {
		t.Errorf("expected 'hourly', got '%s'", Hourly.ToString())
	}
}

func TestMapTaskRepeatRuleToTaskRepeatRule(t *testing.T) {
	// Verify mapping of string to repeat rule
	tests := []struct {
		input    string
		expected TaskRepeatRule
	}{
		{"no_repeat", NoRepeat},
		{"daily", Daily},
		{"hourly", Hourly},
		{"monthly", NoRepeat},
		{"", NoRepeat},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			res := MapTaskRepeatRuleToTaskRepeatRule(tc.input)
			if res != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, res)
			}
		})
	}
}
