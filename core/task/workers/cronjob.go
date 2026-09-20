package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/task"
	task_manager "github.com/smtdfc/nagare/core/task/manager"
)

type CronJobWorker struct {
	mu       sync.Mutex
	timers   map[uuid.UUID]*time.Timer
	logger   *logger.BaseLogger
	taskMgr  *task_manager.TaskManager
	eventBus *event_bus.CoreEventBus
}

func (c *CronJobWorker) Do() {
	ctx := context.Background()
	if err := c.FetchAndScheduleTasks(ctx); err != nil {
		c.logger.Error("Failed in cron job cycle", "err", err)
	}

	go func() {
		ch, unsubscribe := c.eventBus.Subscribe(event_bus.RefreshTaskEvent)
		defer unsubscribe()
		for evt := range ch {
			switch evt.GetEventType() {
			case event_bus.RefreshTaskEvent:
				if err := c.FetchAndScheduleTasks(ctx); err != nil {
					c.logger.Error("Failed in cron job cycle", "err", err)
				}
			}
		}
	}()
}

func (c *CronJobWorker) FetchAndScheduleTasks(ctx context.Context) error {
	upcoming, err := c.taskMgr.GetUpcomingTasks(ctx)
	if err != nil {
		c.logger.Error("Failed to get upcoming task", "err", err)
		return err
	}

	c.logger.Info("Upcoming tasks found", "count", len(upcoming))
	err = c.taskMgr.MarkTasksQueued(ctx, upcoming)
	if err != nil {
		c.logger.Error("Failed to mark tasks queued", "err", err)
		return err
	}

	for _, t := range upcoming {
		if t.TriggerRule.StartTime == nil {
			continue
		}

		tTime := t.TriggerRule.StartTime

		startTime := time.Date(
			tTime.Year(), tTime.Month(), tTime.Day(),
			tTime.Hour(), tTime.Minute(), tTime.Second(),
			tTime.Nanosecond(),
			time.Local,
		)

		duration := max(time.Until(startTime), 0)
		taskID := t.ID

		c.mu.Lock()
		if _, exists := c.timers[taskID]; exists {
			c.mu.Unlock()
			continue
		}

		timer := time.AfterFunc(duration, func() {
			c.ExecuteTask(taskID, t)
		})

		c.timers[taskID] = timer
		c.mu.Unlock()
	}

	return nil
}

func (c *CronJobWorker) ExecuteTask(taskID uuid.UUID, task *task.Task) {
	c.mu.Lock()
	delete(c.timers, taskID)
	c.mu.Unlock()

	c.logger.Info("Executing task", "taskID", taskID.String(), "taskName", task.Name)
	ctx := context.Background()

	c.eventBus.Publish(ctx, event_bus.SendEvent, &event_bus.SendMessageEventPayload{
		RequestID: uuid.New().String(),
		SessionID: task.SessionID.String(),
		Text: fmt.Sprintf(`
			<task>
				<id>%s</id>
				<request>%s</request>
				<note>Add \"[Task: %s] triggered\" before the text response.</note>
				<constraint>You MUST reply using the EXACT same language that the user is currently using in their prompt/request.</constraint>
			</task>
		`, taskID, task.Prompt, taskID),
		SenderType: event_bus.System,
		SenderID:   "",
	})
}

// @Injectable
func NewCronJobWorker(
	logger *logger.BaseLogger,
	taskMgr *task_manager.TaskManager,
	eventBus *event_bus.CoreEventBus,
) *CronJobWorker {
	return &CronJobWorker{
		logger:   logger.With("worker", "core:task:cronjob"),
		taskMgr:  taskMgr,
		eventBus: eventBus,
		timers:   make(map[uuid.UUID]*time.Timer),
	}
}
