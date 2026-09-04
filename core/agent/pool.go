package agent

import (
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool/manager"
)

const NAGARE_AGENT_POOL_SIZE = 10

type AgentPool struct {
	Pool    chan *Agent
	toolMgr *manager.ToolManager
	logger  *logger.BaseLogger
}

func (a *AgentPool) Get() *Agent {
	return <-a.Pool
}

func (a *AgentPool) Put(ag *Agent) *AgentPool {
	a.Pool <- ag
	return a
}

func (a *AgentPool) Seed(size int) *AgentPool {
	for _ = range size {
		a.Put(NewAgent(a.toolMgr, a.logger))
	}

	return a
}

// @Injectable
func NewAgentPool(toolMgr *manager.ToolManager, logger *logger.BaseLogger) *AgentPool {
	return &AgentPool{
		Pool:    make(chan *Agent, NAGARE_AGENT_POOL_SIZE),
		toolMgr: toolMgr,
		logger:  logger,
	}
}
