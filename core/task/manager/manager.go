package manager

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/task"
)

type TaskManager struct {
	taskRepo   *repositories.TaskRepository
	taskMapper *mappers.TaskMapper
	logger     *logger.BaseLogger
}

func (t *TaskManager) Create(ctx context.Context, name string, sessionID string, prompt string, rule *task.TaskTriggerRule) (*task.Task, error) {
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		t.logger.Error("failed to parse session UUID", "sessionID", sessionID, "error", err)
		return nil, custom_errors.ErrCreateTaskFailed
	}
	var nextRunTime *time.Time
	if rule != nil {
		nextRunTime = rule.StartTime
	}

	taskDomain := &task.Task{
		Name:        name,
		Prompt:      prompt,
		SessionID:   sessionUUID,
		TriggerRule: rule,
		Status:      task.Pending,
		NextRunTime: nextRunTime,
	}

	taskEntity, err := t.taskRepo.Create(
		ctx,
		t.taskMapper.ToEntity(taskDomain),
	)
	if err != nil {
		return nil, custom_errors.ErrCreateTaskFailed
	}

	return t.taskMapper.ToDomain(taskEntity), nil
}

func (t *TaskManager) GetUpcomingTasks(ctx context.Context) ([]*task.Task, error) {
	now := time.Now()
	nextHour := now.Add(3 * time.Hour)
	taskEntities, err := t.taskRepo.GetUpcomingScheduledTasks(ctx, now, nextHour)
	if err != nil {
		return nil, custom_errors.ErrGetUpcomingTaskFailed
	}

	return t.taskMapper.ToDomains(taskEntities), nil
}

func (t *TaskManager) MarkTasksQueued(ctx context.Context, tasks []*task.Task) error {
	for _, tsk := range tasks {
		tsk.Status = task.Queued
	}

	err := t.taskRepo.BatchUpdate(ctx, t.taskMapper.ToEntities(tasks))
	if err != nil {
		return custom_errors.ErrMarkTaskQueuedFailed
	}

	return nil
}

func (t *TaskManager) CancelTasksQueued(ctx context.Context, tasks []*task.Task) error {
	for _, tsk := range tasks {
		tsk.Status = task.Pending
	}

	err := t.taskRepo.BatchUpdate(ctx, t.taskMapper.ToEntities(tasks))
	if err != nil {
		return custom_errors.ErrMarkTaskQueuedFailed
	}

	return nil
}

func (t *TaskManager) CancelAllTaskQueued(ctx context.Context) error {
	err := t.taskRepo.CancelAllTaskQueued(ctx)
	if err != nil {
		return custom_errors.ErrCancelTaskQueuedFailed
	}

	return nil
}

// @Injectable
func NewTaskManager(taskRepo *repositories.TaskRepository, logger *logger.BaseLogger) *TaskManager {
	return &TaskManager{
		taskRepo: taskRepo,
		logger:   logger.With("module", "task-manager"),
	}
}
