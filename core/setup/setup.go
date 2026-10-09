package setup

import (
	"context"

	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/chat"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/memory"
	"github.com/smtdfc/nagare/core/plugin/host"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/core/setup/hooks"
	task_manager "github.com/smtdfc/nagare/core/task/manager"
	task_workers "github.com/smtdfc/nagare/core/task/workers"
	"github.com/smtdfc/nagare/pkgs/paths"
)

type CoreSetup struct {
	agentPool         *agent.Pool
	pluginMgr         *manager.PluginManager
	pluginHost        *host.PluginHost
	chatWorker        *chat.Worker
	taskCronJobWorker *task_workers.CronJobWorker
	taskMgr           *task_manager.TaskManager
	vectorMemory      *memory.VectorMemory
	logger            *logger.BaseLogger
}

func (c *CoreSetup) Setup() error {
	ctx := context.Background()

	c.logger.Info("Starting chat worker")
	c.chatWorker.Do()

	c.logger.Info("Starting task worker")
	c.taskCronJobWorker.Do()

	c.logger.Info("Starting plugin host")
	err := c.pluginHost.Start()
	if err != nil {
		return err
	}

	c.logger.Info("Preparing agent pool")
	c.agentPool.Seed(agent.NAGARE_AGENT_POOL_SIZE)

	c.logger.Info("Starting all plugins")
	err = c.pluginMgr.StartAllPlugin(ctx)
	if err != nil {
		return err
	}

	c.logger.Info("Loading vector index")
	err = c.vectorMemory.Load(paths.VectorIndexFile)
	if err != nil {
		c.logger.Warn("Failed to load vector index, creating a new one", "err", err)
		err = c.vectorMemory.Create(1536, 4)
		if err != nil {
			return err
		}
	}
	c.logger.Info("Vector index loaded successfully")

	c.logger.Info("Running hooks")
	for _, fn := range hooks.GetOnStartHooks() {
		err = fn(ctx)
		if err != nil {
			c.logger.Warn("Failed to run hook", "err", err)
		}
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

	c.logger.Info("stopping plugin host")
	err = c.pluginHost.Stop()
	if err != nil {
		c.logger.Error("failed to stop plugin host", "err", err)
	}
	c.logger.Info("plugin host stopped")

	c.logger.Info("syncing vector index")
	err = c.vectorMemory.Sync(paths.VectorIndexFile)
	if err != nil {
		c.logger.Error("failed to sync vector index", "err", err)
		return err
	}
	c.logger.Info("vector index synced successfully")
	c.logger.Info("Running hooks")
	for _, fn := range hooks.GetOnStopHooks() {
		err = fn(ctx)
		if err != nil {
			c.logger.Warn("Failed to run hook", "err", err)
		}
	}

	return nil
}

// @Injectable
func NewCoreSetup(
	agentPool *agent.Pool,
	pluginMgr *manager.PluginManager,
	pluginHost *host.PluginHost,
	chatWorker *chat.Worker,
	taskCronJobWorker *task_workers.CronJobWorker,
	taskMgr *task_manager.TaskManager,
	vectorMemory *memory.VectorMemory,
	logger *logger.BaseLogger,
) *CoreSetup {
	return &CoreSetup{
		agentPool:         agentPool,
		pluginMgr:         pluginMgr,
		pluginHost:        pluginHost,
		chatWorker:        chatWorker,
		taskCronJobWorker: taskCronJobWorker,
		taskMgr:           taskMgr,
		vectorMemory:      vectorMemory,
		logger:            logger.With("module", "system"),
	}
}
