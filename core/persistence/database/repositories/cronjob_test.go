package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupCronJobTestDB(t *testing.T) (*gorm.DB, *entities.Session) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&entities.Session{}, &entities.CronJob{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	session := &entities.Session{Title: "Session 1"}
	if err := db.Create(session).Error; err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	return db, session
}

func TestCronJobRepository_CreateAndGetUpcoming(t *testing.T) {
	db, session := setupCronJobTestDB(t)
	baseLogger, err := logger.NewBaseLogger()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	repo := NewCronJobRepository(db, baseLogger)
	ctx := context.Background()

	now := time.Now()
	upcomingRun := now.Add(10 * time.Minute)
	job := &entities.CronJob{
		Name:        "Morning report",
		SessionID:   session.ID,
		IsActive:    true,
		Prompt:      "Generate report",
		Status:      "pending",
		CronExpr:    "0 9 * * *",
		NextRunTime: &upcomingRun,
	}

	created, err := repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}
	if created.ID == job.ID && created.Name != "Morning report" {
		t.Fatalf("unexpected job name")
	}

	upcoming, err := repo.GetUpcomingCronJobs(now, now.Add(1*time.Hour))
	if err != nil {
		t.Fatalf("failed to get upcoming: %v", err)
	}
	if len(upcoming) != 1 {
		t.Fatalf("expected 1 upcoming job, got %d", len(upcoming))
	}
}

func TestCronJobRepository_CancelAllQueued(t *testing.T) {
	db, session := setupCronJobTestDB(t)
	baseLogger, _ := logger.NewBaseLogger()
	repo := NewCronJobRepository(db, baseLogger)
	ctx := context.Background()

	now := time.Now()
	job := &entities.CronJob{
		Name:        "Queued Job",
		SessionID:   session.ID,
		IsActive:    true,
		Prompt:      "Run task",
		Status:      "queued",
		CronExpr:    "0 * * * *",
		NextRunTime: &now,
	}
	if _, err := repo.Create(ctx, job); err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	if err := repo.CancelAllQueued(ctx); err != nil {
		t.Fatalf("failed to cancel queued: %v", err)
	}

	var found entities.CronJob
	if err := db.First(&found, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("failed to find: %v", err)
	}
	if found.Status != "pending" {
		t.Fatalf("expected status 'pending', got %q", found.Status)
	}
}
