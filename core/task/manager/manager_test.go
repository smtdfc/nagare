package manager

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/task"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testLogger() *logger.BaseLogger {
	return &logger.BaseLogger{Logger: *slog.New(slog.NewTextHandler(io.Discard, nil))}
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	err = db.AutoMigrate(&entities.Task{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func newTestTaskManager(t *testing.T) (*TaskManager, *gorm.DB) {
	t.Helper()
	db := setupTestDB(t)
	repo := repositories.NewTaskRepository(db, testLogger())
	mapper := mappers.NewTaskMapper()
	mgr := NewTaskManager(mapper, repo, testLogger())
	return mgr, db
}

// TestNewTaskManager verifies the constructor returns a non-nil manager.
func TestNewTaskManager(t *testing.T) {
	mgr, _ := newTestTaskManager(t)
	if mgr == nil {
		t.Fatal("expected non-nil TaskManager")
	}
}

// TestTaskManager_Create verifies task creation with valid inputs.
func TestTaskManager_Create(t *testing.T) {
	mgr, _ := newTestTaskManager(t)
	ctx := context.Background()
	sessionID := uuid.New().String()

	now := time.Now().Add(1 * time.Hour)
	rule := &task.TriggerRule{
		By:        task.Scheduled,
		StartTime: &now,
		Repeat:    task.NoRepeat,
	}

	created, err := mgr.Create(ctx, "test-task", sessionID, "do something", rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created == nil {
		t.Fatal("expected non-nil task")
	}
	if created.Name != "test-task" {
		t.Errorf("expected name 'test-task', got %q", created.Name)
	}
	if created.Prompt != "do something" {
		t.Errorf("expected prompt 'do something', got %q", created.Prompt)
	}
	if created.Status != task.Pending {
		t.Errorf("expected status Pending, got %q", created.Status)
	}
}

// TestTaskManager_Create_InvalidSessionID verifies error on invalid UUID.
func TestTaskManager_Create_InvalidSessionID(t *testing.T) {
	mgr, _ := newTestTaskManager(t)
	ctx := context.Background()

	_, err := mgr.Create(ctx, "test-task", "not-a-uuid", "prompt", nil)
	if err == nil {
		t.Fatal("expected error for invalid session UUID")
	}
}

// TestTaskManager_Create_NilTriggerRule verifies task creation with no trigger rule.
func TestTaskManager_Create_NilTriggerRule(t *testing.T) {
	mgr, _ := newTestTaskManager(t)
	ctx := context.Background()
	sessionID := uuid.New().String()

	created, err := mgr.Create(ctx, "simple-task", sessionID, "do it", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.NextRunTime != nil {
		t.Errorf("expected nil NextRunTime, got %v", created.NextRunTime)
	}
}

// TestTaskManager_GetUpcomingTasks verifies retrieval of tasks within the upcoming window.
func TestTaskManager_GetUpcomingTasks(t *testing.T) {
	mgr, db := newTestTaskManager(t)

	// Insert a task with NextRunTime within 3 hours
	nextRun := time.Now().Add(1 * time.Hour)
	taskEntity := &entities.Task{
		ID:          uuid.New(),
		Name:        "upcoming-task",
		SessionID:   uuid.New(),
		IsActive:    true,
		Prompt:      "test prompt",
		Status:      "pending",
		TriggerBy:   "scheduled",
		RepeatType:  "no_repeat",
		NextRunTime: &nextRun,
	}
	db.Create(taskEntity)

	tasks, err := mgr.GetUpcomingTasks()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) == 0 {
		t.Fatal("expected at least 1 upcoming task")
	}
}

// TestTaskManager_MarkTasksQueued verifies tasks are marked as queued.
func TestTaskManager_MarkTasksQueued(t *testing.T) {
	mgr, _ := newTestTaskManager(t)
	ctx := context.Background()
	sessionID := uuid.New().String()

	now := time.Now().Add(1 * time.Hour)
	rule := &task.TriggerRule{
		By:        task.Scheduled,
		StartTime: &now,
		Repeat:    task.NoRepeat,
	}

	created, err := mgr.Create(ctx, "queue-test", sessionID, "prompt", rule)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.MarkTasksQueued(ctx, []*task.Task{created})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.Status != task.Queued {
		t.Errorf("expected status Queued, got %q", created.Status)
	}
}

// TestTaskManager_CancelTasksQueued verifies tasks revert to pending after cancel.
func TestTaskManager_CancelTasksQueued(t *testing.T) {
	mgr, _ := newTestTaskManager(t)
	ctx := context.Background()
	sessionID := uuid.New().String()

	now := time.Now().Add(1 * time.Hour)
	rule := &task.TriggerRule{
		By:        task.Scheduled,
		StartTime: &now,
		Repeat:    task.NoRepeat,
	}

	created, err := mgr.Create(ctx, "cancel-test", sessionID, "prompt", rule)
	if err != nil {
		t.Fatal(err)
	}

	_ = mgr.MarkTasksQueued(ctx, []*task.Task{created})
	err = mgr.CancelTasksQueued(ctx, []*task.Task{created})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.Status != task.Pending {
		t.Errorf("expected status Pending, got %q", created.Status)
	}
}

// TestTaskManager_CancelAllTaskQueued verifies batch cancel of queued tasks.
func TestTaskManager_CancelAllTaskQueued(t *testing.T) {
	mgr, _ := newTestTaskManager(t)
	ctx := context.Background()

	err := mgr.CancelAllTaskQueued(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
