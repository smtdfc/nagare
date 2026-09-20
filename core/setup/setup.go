package setup

import (
	"context"

	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/chat"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/plugin/manager"
	task_manager "github.com/smtdfc/nagare/core/task/manager"
	task_workers "github.com/smtdfc/nagare/core/task/workers"
)

type CoreSetup struct {
	agentPool         *agent.Pool
	pluginMgr         *manager.PluginManager
	chatWorker        *chat.ChatWorker
	taskCronJobWorker *task_workers.CronJobWorker
	taskMgr           *task_manager.TaskManager
	logger            *logger.BaseLogger
}

func (c *CoreSetup) Setup(port string) error {
	ctx := context.Background()

	c.logger.Info("Starting chat worker")
	c.chatWorker.Do()

	c.logger.Info("Starting task worker")
	c.taskCronJobWorker.Do()

	c.logger.Info("Starting plugin host")
	c.pluginMgr.SetPluginHostPort(port)

	c.logger.Info("Preparing agent pool")
	c.agentPool.Seed(agent.NAGARE_AGENT_POOL_SIZE)

	c.logger.Info("Starting all plugins")
	err := c.pluginMgr.StartAllPlugin(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *CoreSetup) Teardown(ctx context.Context) error {
	c.logger.Info("cancelling all tasks")
	err := c.taskMgr.CancelAllTaskQueued(ctx)
	if err != nil {
		c.logger.Error("failed to cancel task(s) queued", "err", err)
		return err
	}
	c.logger.Info("all tasks cancelled")

	c.logger.Info("stopping all plugins")
	err = c.pluginMgr.StopAllPlugin(ctx)
	if err != nil {
		c.logger.Error("failed to stop all plugin", "err", err)
	}
	c.logger.Info("all plugins stopped")

	return nil
}

// @Injectable
func NewCoreSetup(
	agentPool *agent.Pool,
	pluginMgr *manager.PluginManager,
	chatWorker *chat.ChatWorker,
	taskCronJobWorker *task_workers.CronJobWorker,
	taskMgr *task_manager.TaskManager,
	logger *logger.BaseLogger,
) *CoreSetup {
	return &CoreSetup{
		agentPool:         agentPool,
		pluginMgr:         pluginMgr,
		chatWorker:        chatWorker,
		taskCronJobWorker: taskCronJobWorker,
		taskMgr:           taskMgr,
		logger:            logger.With("module", "system"),
	}
}
