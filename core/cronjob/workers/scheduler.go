package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/cronjob"
	cronjob_manager "github.com/smtdfc/nagare/core/cronjob/manager"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type CronJobWorker struct {
	mu         sync.Mutex
	timers     map[uuid.UUID]*time.Timer
	logger     *logger.BaseLogger
	cronJobMgr *cronjob_manager.CronJobManager
	eventBus   *event_bus.CoreEventBus
}

func (c *CronJobWorker) Do() {
	ctx := context.Background()
	if err := c.FetchAndScheduleJobs(ctx); err != nil {
		c.logger.Error("Failed in cronjob cycle", "err", err)
	}

	go func() {
		ch, unsubscribe := c.eventBus.Subscribe(event_bus.RefreshCronJobEvent)
		defer unsubscribe()
		for evt := range ch {
			switch evt.GetEventType() {
			case event_bus.RefreshCronJobEvent:
				if err := c.FetchAndScheduleJobs(ctx); err != nil {
					c.logger.Error("Failed in cronjob refresh cycle", "err", err)
				}
			}
		}
	}()
}

func (c *CronJobWorker) FetchAndScheduleJobs(ctx context.Context) error {
	upcoming, err := c.cronJobMgr.GetUpcomingCronJobs()
	if err != nil {
		c.logger.Error("Failed to get upcoming cronjobs", "err", err)
		return err
	}

	c.logger.Info("Upcoming cronjobs found", "count", len(upcoming))
	err = c.cronJobMgr.MarkCronJobsQueued(ctx, upcoming)
	if err != nil {
		c.logger.Error("Failed to mark cronjobs queued", "err", err)
		return err
	}

	for _, job := range upcoming {
		if job.NextRunTime == nil {
			continue
		}

		duration := max(time.Until(*job.NextRunTime), 0)
		jobID := job.ID

		c.mu.Lock()
		if _, exists := c.timers[jobID]; exists {
			c.mu.Unlock()
			continue
		}

		currentJob := job
		timer := time.AfterFunc(duration, func() {
			c.ExecuteJob(jobID, currentJob)
		})

		c.timers[jobID] = timer
		c.mu.Unlock()
	}

	return nil
}

func (c *CronJobWorker) ExecuteJob(jobID uuid.UUID, job *cronjob.CronJob) {
	c.mu.Lock()
	delete(c.timers, jobID)
	c.mu.Unlock()

	c.logger.Info("Executing cronjob", "jobID", jobID.String(), "jobName", job.Name)
	ctx := context.Background()

	// Advance cronjob next run time & update status
	if err := c.cronJobMgr.AdvanceCronJob(ctx, job); err != nil {
		c.logger.Error("Failed to advance cronjob", "jobID", jobID.String(), "err", err)
	}

	c.eventBus.Publish(ctx, event_bus.SendChatMessageEvent, &event_bus.SendChatMessageEventPayload{
		RequestID: uuid.New().String(),
		SessionID: job.SessionID.String(),
		Messages: messages.ListMessage{
			messages.NewTextMessage(messages.SYSTEM, fmt.Sprintf(`
			<cronjob>
				<id>%s</id>
				<name>%s</name>
				<request>%s</request>
				<note>Add \"[CronJob: %s] triggered\" before the text response.</note>
				<constraint>You MUST reply using the EXACT same language that the user is currently using in their prompt/request.</constraint>
			</cronjob>
		`, jobID, job.Name, job.Prompt, jobID)),
		},
		SenderType: event_bus.System,
		SenderID:   "",
	})

	// Re-check upcoming jobs after advancing
	_ = c.FetchAndScheduleJobs(ctx)
}

// @Injectable
func NewCronJobWorker(
	logger *logger.BaseLogger,
	cronJobMgr *cronjob_manager.CronJobManager,
	eventBus *event_bus.CoreEventBus,
) *CronJobWorker {
	return &CronJobWorker{
		logger:     logger.With("worker", "core:cronjob:worker"),
		cronJobMgr: cronJobMgr,
		eventBus:   eventBus,
		timers:     make(map[uuid.UUID]*time.Timer),
	}
}
