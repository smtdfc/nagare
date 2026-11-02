package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCronJobEntity(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	if err := db.AutoMigrate(&Session{}, &CronJob{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	session := &Session{
		Title: "Test Session",
	}
	if err := db.Create(session).Error; err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	now := time.Now()
	job := &CronJob{
		Name:        "Test CronJob",
		SessionID:   session.ID,
		IsActive:    true,
		Prompt:      "Say hello",
		Status:      "pending",
		CronExpr:    "0 9 * * *",
		NextRunTime: &now,
	}

	if err := db.Create(job).Error; err != nil {
		t.Fatalf("failed to create cronjob: %v", err)
	}

	if job.ID == uuid.Nil {
		t.Fatal("expected non-nil UUID for cronjob")
	}

	var found CronJob
	if err := db.First(&found, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("failed to find cronjob: %v", err)
	}

	if found.Name != "Test CronJob" || found.CronExpr != "0 9 * * *" {
		t.Errorf("unexpected cronjob data: %+v", found)
	}
}
