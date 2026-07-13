package agent

type AgentPool struct {
	Pool chan *Agent
}

func (a *AgentPool) Get() *Agent {
	return <-a.Pool
}

func (a *AgentPool) Put(ag *Agent) {
	a.Pool <- ag
}

func NewAgentPool(size int) *AgentPool {
	return &AgentPool{
		Pool: make(chan *Agent, size),
	}
}
