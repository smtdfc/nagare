package repositories

import (
	"context"
	"time"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"gorm.io/gorm"
)

type CronJobRepository struct {
	db     *gorm.DB
	logger *logger.BaseLogger
}

func (r *CronJobRepository) Create(ctx context.Context, job *entities.CronJob) (*entities.CronJob, error) {
	err := r.db.WithContext(ctx).Create(&job).Error
	if err != nil {
		r.logger.Error("Failed to create cronjob", "err", err)
		return nil, err
	}
	return job, nil
}

// GetUpcomingCronJobs fetches active, pending cronjobs whose NextRunTime falls within [fromTime, toTime].
func (r *CronJobRepository) GetUpcomingCronJobs(fromTime, toTime time.Time) ([]*entities.CronJob, error) {
	var jobs []*entities.CronJob

	err := r.db.Where(
		"is_active = ? AND status = ? AND next_run_time BETWEEN ? AND ?",
		true,
		"pending",
		fromTime,
		toTime,
	).Find(&jobs).Error

	if err != nil {
		r.logger.Error("Failed to get upcoming cronjobs", "err", err)
		return nil, err
	}

	return jobs, nil
}

func (r *CronJobRepository) BatchUpdate(ctx context.Context, jobs []*entities.CronJob) error {
	batchSize := 100

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(jobs); i += batchSize {
			end := min(i+batchSize, len(jobs))
			batch := jobs[i:end]

			for _, job := range batch {
				if err := tx.Model(&entities.CronJob{}).Where("id = ?", job.ID).Updates(job).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		r.logger.Error("Failed to batch update cronjobs", "err", err)
		return err
	}

	return nil
}

func (r *CronJobRepository) CancelAllQueued(ctx context.Context) error {
	res := r.db.WithContext(ctx).Model(&entities.CronJob{}).Where("status = ?", "queued").Update("status", "pending")
	if res.Error != nil {
		r.logger.Error("Failed to cancel all queued cronjobs", "err", res.Error)
		return res.Error
	}
	return nil
}

// @Injectable
func NewCronJobRepository(db *gorm.DB, logger *logger.BaseLogger) *CronJobRepository {
	return &CronJobRepository{
		db:     db,
		logger: logger.With("module", "cronjob-repository"),
	}
}
