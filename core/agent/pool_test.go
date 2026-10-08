package agent

import (
	"testing"
)

func TestPool_PutAndGet(t *testing.T) {
	// Verify Pool channel put and get
	p := &Pool{
		Pool: make(chan *Agent, NAGARE_AGENT_POOL_SIZE),
	}

	ag := &Agent{state: NewAgentState()}
	p.Put(ag)

	retrieved := p.Get()
	if retrieved != ag {
		t.Errorf("expected retrieved agent to match placed agent")
	}
}
