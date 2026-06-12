package agent

import (
	"github.com/smtdfc/nagare/core/messages"
	"github.com/smtdfc/nagare/core/model"
	"github.com/smtdfc/nagare/core/tool"
)

type AgentPool struct {
	Pool  chan *Agent
	Tools tool.ListTool
	Model model.ChatModel
}

func NewAgentPool(size int, model model.ChatModel, tools tool.ListTool) *AgentPool {
	p := &AgentPool{
		Pool:  make(chan *Agent, size),
		Tools: tools,
		Model: model,
	}
	return p
}

func (p *AgentPool) GetOrNew(model model.ChatModel) *Agent {
	select {
	case a := <-p.Pool:
		return a
	default:
		a := NewAgent(p.Model)
		a.Tools = p.Tools
		return a
	}
}

func (p *AgentPool) Put(a *Agent) {
	a.History = messages.ListMessage{
		SYSTEM_PROMPT,
	} // Reset History

	select {
	case p.Pool <- a:
	default:
	}
}
