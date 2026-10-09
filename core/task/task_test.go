package task

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTask_Construction(t *testing.T) {
	// Verify Task struct construction with fields
	id := uuid.New()
	sessID := uuid.New()
	now := time.Now()

	task := Task{
		ID:        id,
		Name:      "Clean Cache",
		SessionID: sessID,
		IsActive:  true,
		Prompt:    "Delete temporary files",
		Status:    Pending,
		TriggerRule: &TriggerRule{
			By:     Scheduled,
			Repeat: Daily,
		},
		NextRunTime: &now,
	}

	if task.ID != id {
		t.Errorf("expected ID %s, got %s", id, task.ID)
	}
	if task.Name != "Clean Cache" {
		t.Errorf("expected name 'Clean Cache', got '%s'", task.Name)
	}
	if task.Status != Pending {
		t.Errorf("expected status Pending, got %s", task.Status)
	}
	if task.NextRunTime == nil || !task.NextRunTime.Equal(now) {
		t.Errorf("expected next run time to match")
	}
}
