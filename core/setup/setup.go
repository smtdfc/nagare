package setup

import (
	"context"

	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/plugin/manager"
)

type CoreSetup struct {
	agentPool *agent.AgentPool
	pluginMgr *manager.PluginManager
}

func (c *CoreSetup) Setup() {
	ctx := context.Background()
	c.agentPool.Seed(agent.NAGARE_AGENT_POOL_SIZE)
	c.pluginMgr.StartAllPlugin(ctx)
}

// @Injectable
func NewCoreSetup(
	agentPool *agent.AgentPool,
	pluginMgr *manager.PluginManager,
) *CoreSetup {
	return &CoreSetup{
		agentPool: agentPool,
		pluginMgr: pluginMgr,
	}
}
