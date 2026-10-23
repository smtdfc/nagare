package agent

import (
	"context"
	"testing"

	"github.com/smtdfc/nagare/core/message"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type testModelAdapter struct{}

func (t testModelAdapter) GetModels(context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (t testModelAdapter) Send(context.Context, string, messages.ListMessage, tool.ListTool) (message.ReadOnlyChannel, error) {
	//TODO implement me
	panic("implement me")
}

func (t testModelAdapter) Embed(context.Context, string, []string) ([][]float32, error) {
	//TODO implement me
	panic("implement me")
}

func TestAgent_Config(t *testing.T) {
	// Verify Agent builder-style methods
	ag := &Agent{
		state: NewAgentState(),
	}

	ag.WithContext(messages.ListMessage{
		messages.NewTextMessage(messages.USER, "initial ctx"),
	})

	if len(ag.state.CurrentMessage) != 1 {
		t.Errorf("expected 1 context message, got %d", len(ag.state.CurrentMessage))
	}

	ag.Reset()
	if len(ag.state.CurrentMessage) != 0 {
		t.Errorf("expected 0 messages after Reset")
	}
}

func TestAgent_WithModelAdapter(t *testing.T) {
	ag := &Agent{}
	ag.WithModelAdapter(&testModelAdapter{})

	if ag.modelAdapter == nil {
		t.Errorf("expected a model adapter")
	}
}

func TestAgent_Reset(t *testing.T) {
	state := NewAgentState()
	state.LoopCounter = 200

	ag := &Agent{
		model:        "test",
		state:        state,
		modelAdapter: &testModelAdapter{},
	}
	ag.Reset()

	if state.LoopCounter != 0 {
		t.Errorf("expected 0 loops, got %d", state.LoopCounter)
	}

	if ag.model != "" {
		t.Errorf("expected empty string")
	}

	if ag.modelAdapter != nil {
		t.Errorf("expected nil adapter")
	}
}

func TestAgent_WithContext(t *testing.T) {
	state := NewAgentState()
	state.LoopCounter = 200

	listMessage := make([]messages.Message, 20)
	for i := range 20 {
		listMessage[i] = messages.NewTextMessage(messages.USER, "initial ctx")
	}

	ag := &Agent{
		state: state,
	}

	ag.WithContext(listMessage)
	if len(ag.state.CurrentMessage) != len(listMessage) {
		t.Errorf("expected %d messages, got %d", len(listMessage), len(ag.state.CurrentMessage))
	}
}

func TestAgent_DumpState(t *testing.T) {
	state := NewAgentState()
	state.LoopCounter = 200
	ag := &Agent{
		state: state,
	}

	dumped := ag.DumpState()
	if dumped.LoopCounter != state.LoopCounter {
		t.Errorf("expected loop counter %d, got %d", state.LoopCounter, dumped.LoopCounter)
	}
}
