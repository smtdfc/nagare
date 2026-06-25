package agent

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/model"
)

type AgentPool struct {
	Pool    chan *Agent
	ToolReg domains.ListTool
	Model   model.ChatModel
}

func NewAgentPool(size int, model model.ChatModel) *AgentPool {
	p := &AgentPool{
		Pool:  make(chan *Agent, size),
		Model: model,
	}
	return p
}

func (p *AgentPool) GetOrNew() *Agent {
	select {
	case a := <-p.Pool:
		return a
	default:
		a := NewAgent(p.Model)
		return a
	}
}

func (p *AgentPool) Put(a *Agent) {
	a.State = NewAgentLoopState()

	select {
	case p.Pool <- a:
	default:
	}
}
