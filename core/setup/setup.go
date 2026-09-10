package setup

import (
	"context"

	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/plugin/manager"
)

type CoreSetup struct {
	agentPool *agent.Pool
	pluginMgr *manager.PluginManager
}

func (c *CoreSetup) Setup() error {
	ctx := context.Background()
	c.agentPool.Seed(agent.NAGARE_AGENT_POOL_SIZE)
	err := c.pluginMgr.StartAllPlugin(ctx)
	if err != nil {
		return err
	}

	return nil
}

// @Injectable
func NewCoreSetup(
	agentPool *agent.Pool,
	pluginMgr *manager.PluginManager,
) *CoreSetup {
	return &CoreSetup{
		agentPool: agentPool,
		pluginMgr: pluginMgr,
	}
}
