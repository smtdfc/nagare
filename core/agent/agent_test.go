package agent

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/messages"
)

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
