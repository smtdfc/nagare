package agent

type AgentPool struct {
	Pool chan *Agent
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
		a.Put(NewAgent("", nil, nil))
	}

	return a
}

func NewAgentPool(size int) *AgentPool {
	return &AgentPool{
		Pool: make(chan *Agent, size),
	}
}
