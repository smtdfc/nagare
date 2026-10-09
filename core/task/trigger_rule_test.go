package task

import (
	"testing"
	"time"
)

func TestTaskTriggerRule_Construction(t *testing.T) {
	// Verify TaskTriggerRule fields
	now := time.Now()
	later := now.Add(time.Hour)

	rule := TriggerRule{
		By:        Scheduled,
		EventName: "on_boot",
		StartTime: &now,
		EndTime:   &later,
		Repeat:    Daily,
	}

	if rule.By != Scheduled {
		t.Errorf("expected By Scheduled, got %s", rule.By)
	}
	if rule.EventName != "on_boot" {
		t.Errorf("expected EventName 'on_boot', got '%s'", rule.EventName)
	}
	if rule.Repeat != Daily {
		t.Errorf("expected Repeat Daily, got %s", rule.Repeat)
	}
	if rule.StartTime == nil || !rule.StartTime.Equal(now) {
		t.Errorf("expected matching StartTime")
	}
	if rule.EndTime == nil || !rule.EndTime.Equal(later) {
		t.Errorf("expected matching EndTime")
	}
}
