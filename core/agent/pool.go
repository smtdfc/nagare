package agent

import (
	"github.com/smtdfc/nagare/core/hil"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool/manager"
)

const NAGARE_AGENT_POOL_SIZE = 10

type Pool struct {
	Pool        chan *Agent
	toolMgr     *manager.ToolManager
	approvalMgr *hil.Manager
	logger      *logger.BaseLogger
}

func (a *Pool) Get() *Agent {
	return <-a.Pool
}

func (a *Pool) Put(ag *Agent) *Pool {
	a.Pool <- ag
	return a
}

func (a *Pool) Seed(size int) *Pool {
	for range size {
		a.Put(NewAgent(a.toolMgr, a.approvalMgr, a.logger))
	}

	return a
}

// @Injectable
func NewAgentPool(toolMgr *manager.ToolManager, approvalMgr *hil.Manager, logger *logger.BaseLogger) *Pool {
	return &Pool{
		Pool:        make(chan *Agent, NAGARE_AGENT_POOL_SIZE),
		toolMgr:     toolMgr,
		approvalMgr: approvalMgr,
		logger:      logger.With("module", "agent-pool"),
	}
}
