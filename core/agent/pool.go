package agent

import (
	"github.com/smtdfc/nagare/core/model"
	"github.com/smtdfc/nagare/core/tool"
)

type AgentPool struct {
	Pool         chan *Agent
	Model        model.ChatModel
	ToolRegistry *tool.ToolRegistry
}

func NewAgentPool(size int, model model.ChatModel, toolReg *tool.ToolRegistry) *AgentPool {
	p := &AgentPool{
		Pool:         make(chan *Agent, size),
		Model:        model,
		ToolRegistry: toolReg,
	}
	return p
}

func (p *AgentPool) GetOrNew() *Agent {
	select {
	case a := <-p.Pool:
		return a
	default:
		a := NewAgent(p.Model, p.ToolRegistry)
		return a
	}
}

func (p *AgentPool) Put(a *Agent) {
	a.State = NewAgentLoopState(p.ToolRegistry)

	select {
	case p.Pool <- a:
	default:
	}
}
