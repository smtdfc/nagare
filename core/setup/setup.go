package setup

import (
	"context"

	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/chat"
	cronjob_manager "github.com/smtdfc/nagare/core/cronjob/manager"
	cronjob_workers "github.com/smtdfc/nagare/core/cronjob/workers"
	"github.com/smtdfc/nagare/core/external_module"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/plugin/host"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/core/setup/hooks"
)

type CoreSetup struct {
	agentPool          *agent.Pool
	pluginMgr          *manager.PluginManager
	pluginHost         *host.PluginHost
	externalModuleHost *external_module.ExternalModuleHost
	chatWorker         *chat.Worker
	cronJobWorker      *cronjob_workers.CronJobWorker
	cronJobMgr         *cronjob_manager.CronJobManager
	logger             *logger.BaseLogger
}

func (c *CoreSetup) Setup() error {
	ctx := context.Background()

	c.logger.Info("Starting chat worker")
	c.chatWorker.Do()

	c.logger.Info("Starting cronjob worker")
	c.cronJobWorker.Do()

	c.logger.Info("Starting external modules")
	err := c.externalModuleHost.StartAllModules()
	if err != nil {
		return err
	}

	c.logger.Info("Starting plugin host")
	err = c.pluginHost.Start()
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
	c.logger.Info("cancelling all queued cronjobs")
	err := c.cronJobMgr.CancelAllQueued(ctx)
	if err != nil {
		c.logger.Error("failed to cancel queued cronjob(s)", "err", err)
		return err
	}
	c.logger.Info("all queued cronjobs cancelled")

	c.logger.Info("Stopping external modules")
	err = c.externalModuleHost.StopAllModules()
	if err != nil {
		c.logger.Error("failed to stop external modules", "err", err)
	}

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
	cronJobWorker *cronjob_workers.CronJobWorker,
	cronJobMgr *cronjob_manager.CronJobManager,
	externalModuleHost *external_module.ExternalModuleHost,
	logger *logger.BaseLogger,
) *CoreSetup {
	return &CoreSetup{
		agentPool:          agentPool,
		pluginMgr:          pluginMgr,
		pluginHost:         pluginHost,
		chatWorker:         chatWorker,
		cronJobWorker:      cronJobWorker,
		cronJobMgr:         cronJobMgr,
		externalModuleHost: externalModuleHost,
		logger:             logger.With("module", "system"),
	}
}
