package messages

import (
	"testing"
)

func TestAgentMessages(t *testing.T) {
	// Test AgentStartedMessage
	startMsg := NewAgentStartedMessage()
	if startMsg == nil || startMsg.GetMessageType() != AgentStartedMessageType {
		t.Fatalf("unexpected start message: %+v", startMsg)
	}
	startMsg.SetInvokeID("inv-start")
	if startMsg.GetInvokeID() != "inv-start" {
		t.Fatalf("expected invoke ID 'inv-start', got %q", startMsg.GetInvokeID())
	}

	// Test AgentCompletedMessage
	completedMsg := NewAgentCompletedMessage(true, false, 1.25)
	if completedMsg == nil || completedMsg.GetMessageType() != AgentCompletedMessageType {
		t.Fatalf("unexpected completed message: %+v", completedMsg)
	}
	if !completedMsg.Success || completedMsg.Cancel || completedMsg.Duration != 1.25 {
		t.Fatalf("fields mismatch in completed message: %+v", completedMsg)
	}

	// Test AgentErrorMessage
	errMsg := NewAgentErrorMessage("execution failed", "ERR_500")
	if errMsg == nil || errMsg.GetMessageType() != AgentErrorMessageType {
		t.Fatalf("unexpected error message: %+v", errMsg)
	}
	if errMsg.Error != "execution failed" || errMsg.Code != "ERR_500" {
		t.Fatalf("fields mismatch in error message: %+v", errMsg)
	}
}
