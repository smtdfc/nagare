package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
)

func TestTaskRepository_CRUD(t *testing.T) {
	db, l := setupTestDB(t)
	repo := &TaskRepository{db: db, logger: l}
	ctx := context.Background()

	id := uuid.New()
	sessID := uuid.New()
	now := time.Now()

	tsk := &entities.Task{
		ID:          id,
		Name:        "Test Task",
		SessionID:   sessID,
		IsActive:    true,
		Status:      "pending",
		TriggerBy:   "scheduled",
		NextRunTime: &now,
	}

	created, err := repo.Create(ctx, tsk)
	if err != nil {
		t.Fatalf("unexpected error creating task: %v", err)
	}
	if created.ID != id {
		t.Errorf("expected ID %s, got %s", id, created.ID)
	}

	// GetUpcomingScheduledTasks
	fromTime := now.Add(-time.Hour)
	toTime := now.Add(time.Hour)
	tasks, err := repo.GetUpcomingScheduledTasks(fromTime, toTime)
	if err != nil {
		t.Fatalf("unexpected error getting upcoming tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 upcoming task, got %d", len(tasks))
	}
}
