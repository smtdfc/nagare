package repositories

import (
	"context"
	"time"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db     *gorm.DB
	logger *logger.BaseLogger
}

func (t *TaskRepository) Create(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	err := t.db.WithContext(ctx).Create(&task).Error
	if err != nil {
		t.logger.Error("Failed to create task", "err", err)
		return nil, err
	}

	return task, nil
}

func (t *TaskRepository) GetUpcomingScheduledTasks(ctx context.Context, fromTime, toTime time.Time) ([]*entities.Task, error) {
	var tasks []*entities.Task

	err := t.db.Where(
		"trigger_by = ? AND is_active = ? AND status = ? AND next_run_time BETWEEN ? AND ? AND (end_time IS NULL OR end_time >= ?)",
		"scheduled",
		true,
		"pending",
		fromTime,
		toTime,
		fromTime,
	).Find(&tasks).Error

	if err != nil {
		t.logger.Error("Failed to get upcoming scheduled tasks", "err", err)
		return nil, err
	}

	return tasks, nil
}

func (t *TaskRepository) BatchUpdate(ctx context.Context, tasks []*entities.Task) error {
	batchSize := 100

	err := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(tasks); i += batchSize {
			end := min(i+batchSize, len(tasks))
			batch := tasks[i:end]

			for _, task := range batch {
				if err := tx.Model(&entities.Task{}).Where("id = ?", task.ID).Updates(task).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		t.logger.Error("Failed to batch update", "err", err)
		return err
	}

	return nil
}

func (t *TaskRepository) CancelAllTaskQueued(ctx context.Context) error {
	res := t.db.WithContext(ctx).Model(&entities.Task{}).Where("status = ?", "queued").Update("status", "pending")
	if res.Error != nil {
		t.logger.Error("Failed to cancel all task queued", "err", res.Error)
		return res.Error
	}
	return nil
}

// @Injectable
func NewTaskRepository(db *gorm.DB, logger *logger.BaseLogger) *TaskRepository {
	return &TaskRepository{
		db:     db,
		logger: logger.With("module", "task-repository"),
	}
}
