//go:build !goverter

package manager

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/smtdfc/nagare/core/cronjob"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	cronjob_mapper "github.com/smtdfc/nagare/core/mappers/cronjob"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
)

var cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

type CronJobManager struct {
	cronJobRepo   *repositories.CronJobRepository
	cronJobMapper cronjob_mapper.ICronJobMapper
	logger        *logger.BaseLogger
}

// Create creates a new CronJob with the given cron expression.
// It parses the expression immediately to compute NextRunTime.
func (m *CronJobManager) Create(ctx context.Context, name string, sessionID string, prompt string, cronExpr string) (*cronjob.CronJob, error) {
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		m.logger.Error("failed to parse session UUID", "sessionID", sessionID, "error", err)
		return nil, custom_errors.ErrCreateCronJobFailed
	}

	expr, err := cronParser.Parse(cronExpr)
	if err != nil {
		m.logger.Error("failed to parse cron expression", "cronExpr", cronExpr, "error", err)
		return nil, custom_errors.ErrCreateCronJobFailed
	}

	nextRun := expr.Next(time.Now())
	job := &cronjob.CronJob{
		Name:        name,
		Prompt:      prompt,
		SessionID:   sessionUUID,
		Status:      cronjob.Pending,
		CronExpr:    cronExpr,
		NextRunTime: &nextRun,
	}

	entity, err := m.cronJobRepo.Create(ctx, m.cronJobMapper.ToEntity(job))
	if err != nil {
		return nil, custom_errors.ErrCreateCronJobFailed
	}

	return m.cronJobMapper.ToDomain(entity), nil
}

func (m *CronJobManager) GetUpcomingCronJobs() ([]*cronjob.CronJob, error) {
	now := time.Now()
	nextHorizon := now.Add(3 * time.Hour)
	entities, err := m.cronJobRepo.GetUpcomingCronJobs(now, nextHorizon)
	if err != nil {
		return nil, custom_errors.ErrGetUpcomingCronJobFailed
	}
	return m.cronJobMapper.ToDomains(entities), nil
}

func (m *CronJobManager) MarkCronJobsQueued(ctx context.Context, jobs []*cronjob.CronJob) error {
	for _, j := range jobs {
		j.Status = cronjob.Queued
	}
	if err := m.cronJobRepo.BatchUpdate(ctx, m.cronJobMapper.ToEntities(jobs)); err != nil {
		return custom_errors.ErrMarkCronJobQueuedFailed
	}
	return nil
}

func (m *CronJobManager) AdvanceCronJob(ctx context.Context, job *cronjob.CronJob) error {
	expr, err := cronParser.Parse(job.CronExpr)
	if err != nil {
		m.logger.Error("failed to parse cron expression when advancing", "cronExpr", job.CronExpr, "error", err)
		return custom_errors.ErrAdvanceCronJobFailed
	}

	now := time.Now()
	nextRun := expr.Next(now)
	job.NextRunTime = &nextRun
	job.LastRunTime = &now
	job.Status = cronjob.Pending

	if err := m.cronJobRepo.BatchUpdate(ctx, m.cronJobMapper.ToEntities([]*cronjob.CronJob{job})); err != nil {
		return custom_errors.ErrAdvanceCronJobFailed
	}

	return nil
}

func (m *CronJobManager) CancelAllQueued(ctx context.Context) error {
	if err := m.cronJobRepo.CancelAllQueued(ctx); err != nil {
		return custom_errors.ErrCancelCronJobQueuedFailed
	}
	return nil
}

// @Injectable
func NewCronJobManager(cronJobMapper cronjob_mapper.ICronJobMapper, cronJobRepo *repositories.CronJobRepository, logger *logger.BaseLogger) *CronJobManager {
	return &CronJobManager{
		cronJobRepo:   cronJobRepo,
		cronJobMapper: cronJobMapper,
		logger:        logger.With("module", "cronjob-manager"),
	}
}
