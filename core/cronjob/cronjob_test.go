package cronjob

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCronJobStruct(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	sessionID := uuid.New()

	job := &CronJob{
		ID:          id,
		Name:        "Test Job",
		SessionID:   sessionID,
		IsActive:    true,
		Prompt:      "Do something",
		Status:      Pending,
		CronExpr:    "0 9 * * *",
		NextRunTime: &now,
	}

	if job.ID != id {
		t.Errorf("ID mismatch")
	}
	if job.CronExpr != "0 9 * * *" {
		t.Errorf("CronExpr mismatch")
	}
	if job.Status != Pending {
		t.Errorf("Status mismatch")
	}
}
