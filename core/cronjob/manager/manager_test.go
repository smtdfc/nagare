package manager

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/cronjob"
	"github.com/smtdfc/nagare/core/logger"
	cronjob_mapper "github.com/smtdfc/nagare/core/mappers/cronjob"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCronJobManager_CreateAndAdvance(t *testing.T) {
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

	baseLogger, _ := logger.NewBaseLogger()
	repo := repositories.NewCronJobRepository(db, baseLogger)
	mapper := cronjob_mapper.NewCronJobMapper()
	mgr := NewCronJobManager(mapper, repo, baseLogger)

	ctx := context.Background()

	// 1. Create with valid cron expression
	job, err := mgr.Create(ctx, "Hourly Ping", session.ID.String(), "Ping server", "0 * * * *")
	if err != nil {
		t.Fatalf("failed to create cronjob: %v", err)
	}
	if job.NextRunTime == nil {
		t.Fatal("expected NextRunTime to be computed")
	}
	if job.Status != cronjob.Pending {
		t.Fatalf("expected Pending status, got %s", job.Status)
	}

	// 2. Advance job
	prevNextRun := *job.NextRunTime
	time.Sleep(10 * time.Millisecond)
	if err := mgr.AdvanceCronJob(ctx, job); err != nil {
		t.Fatalf("failed to advance cronjob: %v", err)
	}
	if job.LastRunTime == nil {
		t.Fatal("expected LastRunTime to be set")
	}
	if job.NextRunTime.Before(prevNextRun) {
		t.Fatalf("expected new NextRunTime to advance")
	}

	// 3. Invalid session ID
	_, err = mgr.Create(ctx, "Bad", "not-a-uuid", "Ping", "0 * * * *")
	if err == nil {
		t.Fatal("expected error on invalid session ID")
	}

	// 4. Invalid cron expression
	_, err = mgr.Create(ctx, "Bad", uuid.New().String(), "Ping", "invalid-cron")
	if err == nil {
		t.Fatal("expected error on invalid cron expr")
	}
}
